package path

import (
	"github.com/text3cn/goodle/container"
)

type Service interface {
	LogPath()
}

type PathService struct {
	Service
	container container.Container
	rootPath  string // 运行时根目录
	// 一些其他路径，比如日志输出路径等，通过配置文件传递进来，配置文件目录等
}

func (s *PathService) LogPath() {}
