// Code generated by microgen 0.9.0. DO NOT EDIT.

package service

import (
	"context"
	log "github.com/go-kit/log"
	service "github.com/recolabs/microgen/examples/generated"
)

// ErrorLoggingMiddleware writes to logger any error, if it is not nil.
func ErrorLoggingMiddleware(logger log.Logger) Middleware {
	return func(next service.StringService) service.StringService {
		return &errorLoggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type errorLoggingMiddleware struct {
	logger log.Logger
	next   service.StringService
}

func (M errorLoggingMiddleware) Uppercase(ctx context.Context, stringsMap map[string]string) (ans string, err error) {
	defer func() {
		if err != nil {
			M.logger.Log("method", "Uppercase", "message", err)
		}
	}()
	return M.next.Uppercase(ctx, stringsMap)
}

func (M errorLoggingMiddleware) Count(ctx context.Context, text string, symbol string) (count int, positions []int, err error) {
	defer func() {
		if err != nil {
			M.logger.Log("method", "Count", "message", err)
		}
	}()
	return M.next.Count(ctx, text, symbol)
}

func (M errorLoggingMiddleware) TestCase(ctx context.Context, comments []*service.Comment) (tree map[string]int, err error) {
	defer func() {
		if err != nil {
			M.logger.Log("method", "TestCase", "message", err)
		}
	}()
	return M.next.TestCase(ctx, comments)
}

func (M errorLoggingMiddleware) DummyMethod(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			M.logger.Log("method", "DummyMethod", "message", err)
		}
	}()
	return M.next.DummyMethod(ctx)
}

func (M errorLoggingMiddleware) IgnoredMethod() {
	M.next.IgnoredMethod()
}

func (M errorLoggingMiddleware) IgnoredErrorMethod() error {
	return M.next.IgnoredErrorMethod()
}
