package session

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (lmw *loggingMiddleware) Online(ctx context.Context, clientId, platform, gateServerIp, gateServerPort string, createdAt uint64) (err error) {
	defer func(begin time.Time) {
		lmw.logger.Log(
			"method", "Online",
			"clientId", clientId,
			"platform", platform,
			"gateServerIp", gateServerIp,
			"gateServerPort", gateServerPort,
			"createdAt", createdAt,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return lmw.next.Online(ctx, clientId, platform, gateServerIp, gateServerPort, createdAt)
}
