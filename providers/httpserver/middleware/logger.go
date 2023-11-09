package middleware

import (
	"fmt"
	"github.com/text3cn/goodle/providers/httpserver"
)

func goodlog() httpserver.MiddlewareHandler {
	return func(c *httpserver.Context) error {
		fmt.Println("Use goodlog middleware")
		c.Next()
		return nil
	}
}
