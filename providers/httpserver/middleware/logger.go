package middleware

import (
	"fmt"
	"github.com/text3cn/t3web/providers/httpserver"
)

func Logger() httpserver.RequestHandler {
	return func(c *httpserver.Context) error {
		fmt.Println("use logger middleware")
		c.Next()
		return nil
	}
}
