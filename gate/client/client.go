package client

import (
	"context"
	"io"
	"pushd/pb"
	"pushd/gate"
	"time"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcd"
	"github.com/go-kit/kit/sd/lb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

func New(etcdAddrs []string, logger log.Logger) (gate.Service, error) {
	factory := factoryFor()
	opts := etcd.ClientOptions{}
	cli, err := etcd.NewClient(context.Background(), etcdAddrs, opts)
	if nil != err {
		return nil, err
	}

	var (
		retryMax     = 3
		retryTimeout = 500 * time.Millisecond
	)
	subscriber, err := etcd.NewSubscriber(cli, "/Gate/Push", factory, logger)
	if nil != err {
		return nil, err
	}
	balancer := lb.NewRoundRobin(subscriber)
	retry := lb.Retry(retryMax, retryTimeout, balancer)

	return &gate.Endpoints{PushEndpoint: retry}, nil
}

func factoryFor() sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			return nil, nil, err
		}
		//		defer conn.Close()
		return grpctransport.NewClient(
			conn,
			"Gate",
			"Push",
			gate.EncodeGRPCOnlineRequest,
			gate.DecodeGRPCOnlineResponse,
			pb.GatePushResponse{},
		).Endpoint(), nil, nil
	}
}
