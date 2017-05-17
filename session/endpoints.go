package session

import (
	"context"
	//	"fmt"

	"github.com/go-kit/kit/endpoint"
	//	"github.com/go-kit/kit/log"
	//	"github.com/go-kit/kit/metrics"
)

type Endpoints struct {
	OnlineEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		OnlineEndpoint: MakeOnlineEndpoint(s),
	}
}

func (e *Endpoints) Online(ctx context.Context, clientId, platform, gateServerIp, gateServerPort string, createdAt uint64) error {
	req := onlineRequest{ClientId: clientId, Platform: platform, GateServerIp: gateServerIp, GateServerPort: gateServerPort, CreatedAt: createdAt}
	resp, err := e.OnlineEndpoint(ctx, req)
	if nil != err {
		return err
	}

	return resp.(onlineResponse).Err
}

func MakeOnlineEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(onlineRequest)
		err = s.Online(ctx, req.ClientId, req.Platform, req.GateServerIp, req.GateServerPort, req.CreatedAt)
		return onlineResponse{
			Err: err,
		}, err
	}
}

type onlineRequest struct {
	ClientId       string
	Platform       string
	GateServerIp   string
	GateServerPort string
	CreatedAt      uint64
}

type onlineResponse struct {
	Err error `json:"err,omitempty"`
}

func (r onlineResponse) error() error { return r.Err }
