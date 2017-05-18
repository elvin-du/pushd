package main

import (
	"os"
	"pushd/pb"

	oldCtx "golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	cli, err := grpc.Dial("127.0.0.1:5502", grpc.WithInsecure())
	if nil != err {
		logger.Log("err", err)
		return
	}
	//	client := pb.NewSessionClient(cli)
	//	_, err = client.Online(oldCtx.Background(), &pb.SessionOnlineRequest{})
	//	if nil != err {
	//		logger.Log("err", err)
	//		return
	//	}
	client := pb.NewGateClient(cli)
	resp, err := client.Push(oldCtx.Background(), &pb.GatePushRequest{
		ClientId: "12345",
		Content:  "hi world",
		Extra:    "",
		Kind:     2,
	})
	if nil != err {
		logger.Log("err", err)
		return
	}
	logger.Log("resp", resp)
}
