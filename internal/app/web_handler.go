package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	Create(ctx *gin.Context)
}

func NewWebHandler(appConfig *config.AppConfig, service AppService) *webHandler {
	return &webHandler{appConfig, service}
}

func (r *webHandler) Default(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{"title": r.AppConfig.Web.Title, "data": "{'content':''}"})
}

func (r *webHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	a, _ := r.Service.GetbyId(id)
	x, _ := json.Marshal(a[0])
	c.HTML(http.StatusOK, "index", gin.H{"id": id, "data": string(x), "removable": true})
}

func (r *webHandler) DeleteById(c *gin.Context) {
	id := c.Param("id")
	formData := &model.FormData{}
	c.Bind(formData)

	if errs := r.Service.DeleteById(id, formData); errs != nil {
		c.Writer.WriteHeader(http.StatusOK)
		var errElement strings.Builder
		for _, err := range errs {
			errElement.WriteString(fmt.Sprintf("<li>%s</li>", err))
		}
		c.Writer.Write([]byte("<ul>" + errElement.String() + "</ul>"))
	} else {
		c.Header("HX-Redirect", "/"+formData.Id)
	}
}

func (r *webHandler) Save(c *gin.Context) {
	id := c.Param("id")
	formData := &model.FormData{}
	c.Bind(formData)

	if err := r.Service.Save(id, formData); err != nil {
		x, _ := json.Marshal(formData)
		c.HTML(http.StatusOK, "index_content", gin.H{"id": id, "data": string(x), "errors": err})
	} else {
		c.Header("HX-Redirect", "/"+formData.Id)
	}
}

func (r *webHandler) Create(c *gin.Context) {
	formData := &model.FormData{}
	c.Bind(formData)

	if err := r.Service.Create(formData); err != nil {
		x, _ := json.Marshal(formData)
		c.HTML(http.StatusOK, "index_content", gin.H{"data": string(x), "errors": err})
	} else {
		c.Header("HX-Redirect", "/"+formData.Id)
	}
}
