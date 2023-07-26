package httpserver

import (
	"github.com/text3cn/t3web/container"
)

type HttpServerService struct {
	container container.Container
	*T3WebEngine
}
