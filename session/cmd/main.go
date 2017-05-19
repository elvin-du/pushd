package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"pushd/pb"
	"pushd/session"
	"strings"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcd"
	"github.com/go-kit/kit/tracing/opentracing"
	"google.golang.org/grpc"
	"sourcegraph.com/sourcegraph/appdash"
	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
)

var (
	grpcAddr    *string
	appdashAddr *string
	etcdAddrs   []string
)

func init() {
	grpcAddr = flag.String("grpc.addr", ":5503", "TCP listen address")
	appdashAddr = flag.String("appdash.addr", ":5507", "Enable Appdash tracing via an Appdash server host:port")
	etcdAddresses := flag.String("etcd.addrs", "http://127.0.0.1:2379", "ETCD V2 servers host:port,host:port")
	flag.Parse()
	etcdAddrs = strings.Split(*etcdAddresses, ",")
}

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s session.Service
	{
		s = session.NewService()
		s = session.LoggingMiddleware(logger)(s)
	}

	var endpoints session.Endpoints
	{
		endpoints = session.MakeServerEndpoints(s)
	}

	tracer := appdashot.NewTracer(appdash.NewRemoteCollector(*appdashAddr))
	endpoints.OnlineEndpoint = opentracing.TraceServer(tracer, "Online")(endpoints.OnlineEndpoint)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		errc <- fmt.Errorf("%s", <-c)
	}()

	// gRPC transport.
	go func() {
		logger := log.With(logger, "transport", "gRPC")

		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}

		srv := session.MakeGRPCServer(endpoints, logger)
		s := grpc.NewServer()
		pb.RegisterSessionServer(s, srv)

		logger.Log("addr", *grpcAddr)
		errc <- s.Serve(ln)
	}()

	cli, err := etcd.NewClient(context.Background(), etcdAddrs, etcd.ClientOptions{})
	if nil != err {
		logger.Log("err", err)
		errc <- err
	}
	onlineRegistar := etcd.NewRegistrar(cli, etcd.Service{
		Key:   "/Session/Online/127.0.0.1" + *grpcAddr,
		Value: "127.0.0.1" + *grpcAddr,
	}, logger)
	onlineRegistar.Register()
	defer onlineRegistar.Deregister()

	logger.Log("exit", <-errc)
}
