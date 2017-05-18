package gate

import (
	"context"

	//	stdopentracing "github.com/opentracing/opentracing-go"
	oldcontext "golang.org/x/net/context"

	"pushd/pb"

	"github.com/go-kit/kit/log"
	//	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC AddServer.
func MakeGRPCServer(endpoints Endpoints, logger log.Logger) pb.GateServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		push: grpctransport.NewServer(
			endpoints.PushEndpoint,
			DecodeGRPCPushRequest,
			EncodeGRPCPushResponse,
			append(options)...,
		),
	}
}

type grpcServer struct {
	push grpctransport.Handler
}

func (s *grpcServer) Push(ctx oldcontext.Context, req *pb.GatePushRequest) (*pb.GatePushResponse, error) {
	_, rep, err := s.push.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GatePushResponse), nil
}

func DecodeGRPCPushRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GatePushRequest)
	return pushRequest{ClientId: req.ClientId, Content: req.Content, Extra: req.Extra, Kind: req.Kind}, nil
}

func EncodeGRPCPushResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(pushResponse)
	return &pb.GatePushResponse{}, resp.Err
}

func DecodeGRPCPushResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	return pushResponse{}, nil
}

func EncodeGRPCPushRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(pushRequest)
	return &pb.GatePushRequest{ClientId: req.ClientId, Content: req.Content, Extra: req.Extra, Kind: req.Kind}, nil
}
