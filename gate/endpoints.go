package gate

import (
	"context"
	//	"fmt"

	"github.com/go-kit/kit/endpoint"
	//	"github.com/go-kit/kit/log"
	//	"github.com/go-kit/kit/metrics"
)

type Endpoints struct {
	PushEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PushEndpoint: MakePushEndpoint(s),
	}
}

func (e *Endpoints) Push(ctx context.Context, clientId, content, extra string, kind uint32) error {
	req := pushRequest{ClientId: clientId, Content: content, Extra: extra, Kind: kind}
	resp, err := e.PushEndpoint(ctx, req)
	if nil != err {
		return err
	}

	return resp.(pushResponse).Err
}

func MakePushEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(pushRequest)
		err = s.Push(ctx, req.ClientId, req.Content, req.Extra, req.Kind)
		return pushResponse{
			Err: err,
		}, err
	}
}

type pushRequest struct {
	ClientId string
	Content  string
	Extra    string
	Kind     uint32
}

type pushResponse struct {
	Err error `json:"err,omitempty"`
}

func (r pushResponse) error() error { return r.Err }
