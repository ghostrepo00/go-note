package app

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/ghostrepo00/go-note/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

type webserver struct {
	appConfig *config.AppConfig
}

func NewWebServer(config *config.AppConfig) *webserver {
	return &webserver{config}
}

type test struct {
	Id    int64
	Value string
}

func ConfigureWebRouter(router *gin.Engine, appConfig *config.AppConfig, dbClient *supabase.Client) {
	t := template.Must(template.ParseGlob("web/template/**/*.html"))
	router.SetHTMLTemplate(t)

	router.Static("/assets", "web/assets")
	router.StaticFile("/favicon.ico", "web/favicon.ico")

	router.GET("", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	// ConfigureHomeRouter(router)
	// ConfigureWebtagRouter(router, webtagApp)
}

func (r *webserver) Run() {
	if dbClient, err := supabase.NewClient(r.appConfig.SupabaseUrl, r.appConfig.SupabaseKey, nil); err != nil {
		panic(err)
	} else {
		slog.Info("Database connected")
		if data, _, err := dbClient.From("test").Select("id, value", "exact", false).Execute(); err == nil {
			var result []*test
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
