package app

import (
	"log/slog"
	"net/http"

	"github.com/ghostrepo00/go-note/config"
	"github.com/ghostrepo00/go-note/internal/pkg/model"
	"github.com/gin-gonic/gin"
)

type webHandler struct {
	AppConfig *config.AppConfig
	Service   AppService
}

type WebHandler interface {
	Default(ctx *gin.Context)
	GetById(ctx *gin.Context)
	DeleteById(c *gin.Context)
	Save(ctx *gin.Context)
}

func NewWebHandler(appConfig *config.AppConfig, service AppService) *webHandler {
	return &webHandler{appConfig, service}
}

func (r *webHandler) Default(c *gin.Context) {
	p := make(map[string]string)
	p["title"] = r.AppConfig.Web.Title
	c.HTML(http.StatusOK, "index.html", p)
}

func (r *webHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	slog.Info("request id", "id", id)
	a, _ := r.Service.GetbyId(id)
	data := make(map[string]string)
	data["content"] = a[0].Content
	data["id"] = a[0].Id
	c.HTML(http.StatusOK, "index.html", data)
}

func (r *webHandler) DeleteById(c *gin.Context) {
	id := c.Param("id")
	r.Service.DeletebyId(id)
	slog.Info("request id", "id", id)
	c.Header("HX-Redirect", "/")
}

func (r *webHandler) Save(c *gin.Context) {
	formData := &model.FormData{}
	c.Bind(formData)
	if err := r.Service.Save(formData); err != nil {
		slog.Error("error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.Header("HX-Redirect", "/"+formData.Id)
	}
	slog.Info("request id", "id", formData.Id, "content", formData.Content)
}
