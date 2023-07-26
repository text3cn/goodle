package middleware

import (
	"fmt"
	"github.com/text3cn/goodle/providers/httpserver"
)

func Logger() httpserver.RequestHandler {
	return func(c *httpserver.Context) error {
		fmt.Println("use logger middleware")
		c.Next()
		return nil
	}
}
