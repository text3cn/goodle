package httpserver

import (
	"github.com/text3cn/goodle/container"
)

type HttpServerService struct {
	container container.Container
	*GoodleEngine
}
