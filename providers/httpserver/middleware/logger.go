package middleware

import (
	"fmt"
	"github.com/text3cn/goodle/providers/httpserver"
)

func Logger() httpserver.MiddlewareHandler {
	return func(c *httpserver.Context) error {
		fmt.Println("Use logger middleware")
		c.Next()
		return nil
	}
}
