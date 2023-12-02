package app

import (
	"encoding/json"
	"html/template"
	"log/slog"

	"github.com/ghostrepo00/go-note/config"
	"github.com/ghostrepo00/go-note/internal/pkg/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type webServer struct {
	appConfig *config.AppConfig
}

func NewWebServer(config *config.AppConfig) *webServer {
	return &webServer{config}
}

func ConfigureWebRouter(router *gin.Engine, appConfig *config.AppConfig, dbClient *supabase.Client) {
	t := template.Must(template.ParseGlob("web/template/**/*.html"))
	router.SetHTMLTemplate(t)

	router.Static("/assets", "web/assets")
	router.StaticFile("/favicon.ico", "web/favicon.ico")

	var service AppService = NewAppService(appConfig, dbClient)
	var handler WebHandler = NewWebHandler(appConfig, service)

	router.GET("", handler.Default)
	router.GET("/:id", handler.GetById)
	router.DELETE("/:id", handler.DeleteById)
	router.POST("/save", handler.Save)
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
		router := gin.Default()
		router.Use(cors.Default())
		ConfigureWebRouter(router, r.appConfig, dbClient)
		router.Run(r.appConfig.Web.Host)
	}
}
