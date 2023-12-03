package app

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

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

func ConfigureWebRouter(appConfig *config.AppConfig, dbClient *supabase.Client) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		c.HTML(http.StatusNotFound, "error", gin.H{"Status": 500, "Message": "Not Found", "Description": err.(error)})
	}))

	t := template.Must(template.ParseGlob("web/template/**/*.html"))
	router.SetHTMLTemplate(t)

	router.Static("/assets", "web/assets")
	router.StaticFile("/favicon.ico", "web/favicon.ico")

	var service AppService = NewAppService(appConfig, dbClient)
	var handler WebHandler = NewWebHandler(appConfig, service)

	router.GET("", handler.Default)
	router.GET("/:id", handler.GetById)
	router.DELETE("/:id", handler.DeleteById)
	router.POST("", handler.Save)
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
