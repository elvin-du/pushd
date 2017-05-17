package client

import (
	"context"
	"io"
	"pushd/session"
	//	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcd"
	"github.com/go-kit/kit/sd/lb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

func NewClient() {
	etcdAddrs := []string{""}
	cli, err := etcd.NewClient(context.Background(), etcdAddrs, etcd.ClientOptions{})
	if nil != err {
		return
	}

	conn, err := grpc.Dial(":5504", grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	grpctransport.NewClient()
}

//func New(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) session.Service {
//	limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))

//	var onlineEndpoint endpoint.Endpoint
//	{
//		onlineEndpoint = grpctransport.NewClient(
//			conn,
//			"session",
//			"Online",
//			session.EncodeGRPCOnlineRequest,
//			session.DecodeGRPCOnlineResponse,
//			pb.SessionOnlineResponse{},
//			//			grpctransport.ClientBefore(opentracing.ToGRPCRequest(tracer, logger)),
//		).Endpoint()
//		//		onlineEndpoint = opentracing.TraceClient(tracer, "Sum")(onlineEndpoint)
//		onlineEndpoint = limiter(onlineEndpoint)
//		//		onlineEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
//		//			Name:    "Sum",
//		//			Timeout: 30 * time.Second,
//		//		}))(onlineEndpoint)
//	}

//	return session.Endpoints{
//		OnlineEndpoint: onlineEndpoint,
//	}
//}

//func factoryFor(makeEndpoint func(session.Service) endpoint.Endpoint) sd.Factory {
//	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
//		service, err := session.MakeClientEndpoints(instance)
//		if err != nil {
//			return nil, nil, err
//		}
//		return makeEndpoint(service), nil, nil
//	}
//}

//func MakeClientEndpoints(instance string) (session.Endpoints, error) {
//}
