package httpserver

import (
	"github.com/text3cn/goodle/core"
)

type HttpServerService struct {
	container core.Container
	*Engine
}
