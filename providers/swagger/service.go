package swagger

import (
	"github.com/spf13/cast"
	httpSwagger "github.com/swaggo/http-swagger"
	goodleconfig "github.com/text3cn/goodle/config"
	"github.com/text3cn/goodle/container"
	"net/http"
)

type Service interface {
	Init()
}

type SwaggerService struct {
	Service
	c container.Container
}

func (self *SwaggerService) Init() {
	configService := goodleconfig.Instance()
	cfg := configService.GetSwagger()
	port := cast.ToString(cfg.SwaggerUiPort)
	http.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.FilePath)
	})
	http.HandleFunc("/", httpSwagger.Handler(
		httpSwagger.URL("http://"+cfg.SwaggerUiHost+":"+port+"/swagger"),
	))
	go http.ListenAndServe(":"+port, nil)
}
