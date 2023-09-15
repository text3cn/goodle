package engine

import (
	"github.com/text3cn/goodle/container"
	"github.com/text3cn/goodle/providers/cache"
	"github.com/text3cn/goodle/providers/httpserver"
	"github.com/text3cn/goodle/providers/orm"
	"github.com/text3cn/goodle/providers/redis"
	"github.com/text3cn/goodle/providers/swagger"
)

func bindServices(c container.Container) {
	c.Bind(&httpserver.HttpServerProvider{})
	c.Bind(&cache.CacheServiceProvider{})
	c.Bind(&orm.OrmProvider{})
	c.Bind(&redis.ReidsProvider{})
	c.Bind(&swagger.SwaggerProvider{})
}
