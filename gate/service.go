package gate

import (
	"context"
	"fmt"
	"pushd/session/client"

	"os"

	"github.com/go-kit/kit/log"
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
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	svc, err := client.New([]string{"http://127.0.0.1:2379"}, "127.0.0.1:5507", logger)
	if nil != err {
		logger.Log("err:", err)
		return err
	}
	err = svc.Online(ctx, clientId, "android", "123", "456", 123456)
	if nil != err {
		logger.Log("err:", err)
		return err
	}
	return nil
}
