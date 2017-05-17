package session

import (
	"context"
	"fmt"
)

type Service interface {
	//	Push(ctx context.Context, clientId, content, extra string, packetId, kind uint32) error
	Online(ctx context.Context, clientId, platform, gateServerIp, gateServerPort string, createdAt uint64) error
}

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) Online(ctx context.Context, clientId, platform, gateServerIp, gateServerPort string, createdAt uint64) error {
	fmt.Printf("%s online on server(%s:%s) at:%d via %s SUCCESS", clientId, gateServerIp, gateServerPort, createdAt, platform)
	return nil
}
