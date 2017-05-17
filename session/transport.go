package session

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
func MakeGRPCServer(endpoints Endpoints, logger log.Logger) pb.SessionServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		online: grpctransport.NewServer(
			endpoints.OnlineEndpoint,
			DecodeGRPCOnlineRequest,
			EncodeGRPCOnlineResponse,
			append(options)...,
		),
	}
}

type grpcServer struct {
	online grpctransport.Handler
}

func (s *grpcServer) Online(ctx oldcontext.Context, req *pb.SessionOnlineRequest) (*pb.SessionOnlineResponse, error) {
	_, rep, err := s.online.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SessionOnlineResponse), nil
}

func DecodeGRPCOnlineRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SessionOnlineRequest)
	return onlineRequest{ClientId: req.ClientId, GateServerIp: req.GateServerIP, GateServerPort: req.GateServerPort, Platform: req.Platform, CreatedAt: req.CreatedAt}, nil
}

func DecodeGRPCOnlineResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	return onlineResponse{}, nil
}

func EncodeGRPCOnlineRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(onlineRequest)
	return &pb.SessionOnlineRequest{
		ClientId:       req.ClientId,
		GateServerIP:   req.GateServerIp,
		GateServerPort: req.GateServerPort,
		Platform:       req.Platform,
		CreatedAt:      req.CreatedAt,
	}, nil
}

func EncodeGRPCOnlineResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(onlineResponse)
	return &pb.SessionOnlineResponse{}, resp.Err
}
