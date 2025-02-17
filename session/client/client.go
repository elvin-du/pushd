package client

import (
	"context"
	"io"
	"pushd/pb"
	"pushd/session"
	"time"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"sourcegraph.com/sourcegraph/appdash"
	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
)

func New(etcdAddrs []string, appdashAddr string, logger log.Logger) (session.Service, error) {
	tracer := appdashot.NewTracer(appdash.NewRemoteCollector(appdashAddr))
	factory := factoryFor(tracer, "Online")
	opts := etcd.ClientOptions{}
	cli, err := etcd.NewClient(context.Background(), etcdAddrs, opts)
	if nil != err {
		return nil, err
	}

	var (
		retryMax     = 3
		retryTimeout = 500 * time.Millisecond
	)
	subscriber, err := etcd.NewSubscriber(cli, "/Session/Online", factory, logger)
	if nil != err {
		return nil, err
	}
	balancer := lb.NewRoundRobin(subscriber)
	retry := lb.Retry(retryMax, retryTimeout, balancer)

	return &session.Endpoints{OnlineEndpoint: retry}, nil
}

func factoryFor(tracer stdopentracing.Tracer, svcName string) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			return nil, nil, err
		}
		//		defer conn.Close()
		endpoint := grpctransport.NewClient(
			conn,
			"Session",
			"Online",
			session.EncodeGRPCOnlineRequest,
			session.DecodeGRPCOnlineResponse,
			pb.SessionOnlineResponse{},
		).Endpoint()
		endpoint = opentracing.TraceClient(tracer, svcName)(endpoint)
		return endpoint, nil, nil
	}
}
