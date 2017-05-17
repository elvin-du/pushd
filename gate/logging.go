package gate

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

func (lmw *loggingMiddleware) Push(ctx context.Context, clientId, content, extra string, kind uint32) (err error) {
	defer func(begin time.Time) {
		lmw.logger.Log(
			"method", "push",
			"clientId", clientId,
			"content", content,
			"extra", extra,
			"kind", kind,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return lmw.next.Push(ctx, clientId, content, extra, kind)
}
