package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"pushd/session"
	"pushd/pb"
	"syscall"

	"google.golang.org/grpc"
	"github.com/go-kit/kit/log"
)

var (
	tcpAddr     *string
	grpcAddr    *string
	appdashAddr *string
)

func init() {
	grpcAddr = flag.String("grpc.addr", ":5503", "TCP listen address")
	appdashAddr = flag.String("appdash.addr", "", "Enable Appdash tracing via an Appdash server host:port")
}

func main() {
	flag.Parse()
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

	logger.Log("exit", <-errc)
}
