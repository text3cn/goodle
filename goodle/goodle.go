package goodle

import (
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/cache"
	"github.com/text3cn/goodle/providers/config"
	"github.com/text3cn/goodle/providers/httpserver"
	"github.com/text3cn/goodle/providers/orm"
	"github.com/text3cn/goodle/providers/redis"
	"github.com/text3cn/goodle/providers/swagger"
)

type Goodle struct {
}

func Init(onBoot ...func(container core.Container)) *Goodle {
	container := core.GlobalCore()
	container.Bind(&config.ConfigProvider{})
	container.Bind(&httpserver.HttpServerProvider{})
	container.Bind(&cache.CacheServiceProvider{})
	container.Bind(&orm.OrmProvider{})
	container.Bind(&redis.ReidsProvider{})
	container.Bind(&swagger.SwaggerProvider{})
	if len(onBoot) > 0 {
		onBoot[0](container)
	}
	return &Goodle{}
}
