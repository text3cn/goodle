package swagger

import (
	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/text3cn/goodle/core"
	"github.com/text3cn/goodle/providers/config"
	"net/http"
)

type Service interface {
	Init()
}

type SwaggerService struct {
	Service
	c core.Container
}

func (self *SwaggerService) Init() {
	configService := self.c.NewSingle(core.Config).(config.Service)
	cfg := configService.GetSwagger()
	port := cast.ToString(cfg.SwaggerUiPort)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://"+cfg.SwaggerUiHost+":"+port+"/swagger"),
	))

	go http.ListenAndServe(":"+port, r)

	//http.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, cfg.FilePath)
	//})
	//http.HandleFunc("/s", httpSwagger.Handler(
	//	httpSwagger.URL("http://"+cfg.SwaggerUiHost+":"+port+"/swagger"),
	//))
	//go http.ListenAndServe(":"+port, nil)
}
