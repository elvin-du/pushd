package gate

import (
	"context"
	"fmt"
)

type Service interface {
	//	Push(ctx context.Context, clientId, content, extra string, packetId, kind uint32) error
	Push(ctx context.Context, clientId, content, extra string, kind uint32) error
}

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) Push(ctx context.Context, clientId, content, extra string, kind uint32) error {
	fmt.Printf("push %s to %s SUCCESS", content, clientId)
	return nil
}
