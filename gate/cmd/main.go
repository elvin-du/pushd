package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"pushd/gate"
	"pushd/pb"
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
	tcpAddr     *string
	grpcAddr    *string
	appdashAddr *string
	etcdAddrs   []string
)

func init() {
	tcpAddr = flag.String("tcp.addr", ":5501", "GRPC listen address")
	grpcAddr = flag.String("grpc.addr", ":5502", "TCP listen address")
	appdashAddr = flag.String("appdash.addr", ":5507", "Enable Appdash tracing via an Appdash server host:port")
	etcdAddresses := flag.String("etcd.addrs", "http://127.0.0.1:2379", "ETCD V2 servers host:port,host:port")
	flag.Parse()
	etcdAddrs = strings.Split(*etcdAddresses, ",")
}

func main() {
	flag.Parse()
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s gate.Service
	{
		s = gate.NewService()
		s = gate.LoggingMiddleware(logger)(s)
	}

	var endpoints gate.Endpoints
	{
		endpoints = gate.MakeServerEndpoints(s)
	}

	tracer := appdashot.NewTracer(appdash.NewRemoteCollector(*appdashAddr))
	endpoints.PushEndpoint = opentracing.TraceServer(tracer, "Push")(endpoints.PushEndpoint)

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

		srv := gate.MakeGRPCServer(endpoints, logger)
		s := grpc.NewServer()
		pb.RegisterGateServer(s, srv)

		logger.Log("addr", *grpcAddr)
		errc <- s.Serve(ln)
	}()

	cli, err := etcd.NewClient(context.Background(), etcdAddrs, etcd.ClientOptions{})
	if nil != err {
		logger.Log("err", err)
		errc <- err
	}
	onlineRegistar := etcd.NewRegistrar(cli, etcd.Service{
		Key:   "/Gate/Push/127.0.0.1" + *grpcAddr,
		Value: "127.0.0.1" + *grpcAddr,
	}, logger)
	onlineRegistar.Register()
	defer onlineRegistar.Deregister()

	logger.Log("exit", <-errc)
}
