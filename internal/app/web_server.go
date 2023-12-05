package app

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ghostrepo00/go-note/config"
	"github.com/ghostrepo00/go-note/internal/pkg/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type webServer struct {
	appConfig *config.AppConfig
}

func NewWebServer(config *config.AppConfig) *webServer {
	return &webServer{config}
}

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "web/template/shared/base.html", "web/template/home/index.html")
	r.AddFromFiles("error_list", "web/template/shared/error_list.html")
	r.AddFromFiles("error", "web/template/shared/base.html", "web/template/shared/error.html")
	return r
}

func ConfigureWebRouter(appConfig *config.AppConfig, dbClient *supabase.Client) *gin.Engine {
	router := gin.Default()
	router.HTMLRender = createMyRender()
	router.Use(cors.Default())
	router.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		slog.Error("Unhandled exception", "error", err.(error))
		c.HTML(http.StatusInternalServerError, "error", gin.H{"Status": 500, "Message": "Internal Error"})
	}))

	router.Static("/assets", "web/assets")
	router.StaticFile("/favicon.ico", "web/favicon.ico")

	var service AppService = NewAppService(appConfig, dbClient)
	var handler WebHandler = NewWebHandler(appConfig, service)

	router.GET("", handler.Default)
	router.GET("/:id", handler.GetById)
	router.POST("/:id/delete", handler.DeleteById)
	router.POST("", handler.Create)
	router.POST("/:id", handler.Save)

	return router
}

func (r *webServer) Run() {
	if dbClient, err := supabase.NewClient(r.appConfig.SupabaseUrl, r.appConfig.SupabaseKey, nil); err != nil {
		panic(err)
	} else {
		slog.Info("Database connected")
		if data, _, err := dbClient.From("test").Select("id, value", "exact", false).Execute(); err == nil {
			var result []*model.Test
			err := json.Unmarshal(data, &result)
			if err == nil {
				slog.Info(result[0].Value)
			}
		}
		ConfigureWebRouter(r.appConfig, dbClient).Run(r.appConfig.Web.Host)
	}
}
