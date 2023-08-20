package middleware

import (
	"github.com/text3cn/goodle/providers/httpserver"
)

func Auth(fn func(ctx *httpserver.Context) error) httpserver.MiddlewareHandler {
	return func(context *httpserver.Context) error {
		err := fn(context)
		if err == nil {
			context.Next()
			return nil
		}
		return err
	}
}
