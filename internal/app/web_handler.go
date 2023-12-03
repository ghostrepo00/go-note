package app

import (
	"encoding/json"
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
	c.HTML(http.StatusOK, "index.html", gin.H{"title": r.AppConfig.Web.Title, "data": "{'content':'', 'value':'z'}"})
}

func (r *webHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	slog.Info("request id", "id", id)
	a, _ := r.Service.GetbyId(id)
	// if len(a) == 0 {
	// 	a = append(a, &model.FormData{Id: id})
	// }
	x, _ := json.Marshal(a[0])
	c.HTML(http.StatusOK, "index.html", gin.H{"id": id, "data": string(x)})
}

func (r *webHandler) DeleteById(c *gin.Context) {
	id := c.Param("id")
	r.Service.DeletebyId(id)
	slog.Info("request id", "id", id)
	c.Header("HX-Redirect", "/")
}

func (r *webHandler) Save(c *gin.Context) {
	id := c.Param("id")
	formData := &model.FormData{}
	c.Bind(formData)

	if outputPassword, err := r.Service.Save(id, formData); err != nil {
		slog.Error("error", err, outputPassword)
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.Header("HX-Redirect", "/"+formData.Id)
	}
	slog.Info("request id", "id", formData.Id, "content", formData.Content)
}
