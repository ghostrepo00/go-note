package app

import (
	"encoding/json"
	"errors"
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
	if entity, _ := r.Service.GetbyId(id); len(entity) == 0 {
		c.HTML(http.StatusNotFound, "error", gin.H{"Status": 404, "Message": "Record not found"})
	} else {
		data, _ := json.Marshal(entity[0])
		c.HTML(http.StatusOK, "index", gin.H{"id": id, "data": string(data), "removable": true})
	}
}

func (r *webHandler) DeleteById(c *gin.Context) {
	id := c.Param("id")
	formData := &model.FormData{}
	c.Bind(formData)

	if formData.Password == "" {
		c.HTML(http.StatusOK, "error_list", gin.H{"errors": &[]error{errors.New("Password is required")}})
		return
	}

	if errs := r.Service.DeleteById(id, formData); errs != nil {
		c.HTML(http.StatusOK, "error_list", gin.H{"errors": errs})
	} else {
		c.Header("HX-Redirect", "/")
	}
}

func (r *webHandler) Save(c *gin.Context) {
	id := c.Param("id")
	formData := &model.FormData{}
	c.Bind(formData)

	if formData.Password == "" {
		c.HTML(http.StatusOK, "error_list", gin.H{"errors": &[]error{errors.New("Password is required")}})
		return
	}

	if errs := r.Service.Save(id, formData); errs != nil {
		c.HTML(http.StatusOK, "error_list", gin.H{"errors": errs})
	} else {
		c.Header("HX-Redirect", "/"+formData.Id)
	}
}

func (r *webHandler) Create(c *gin.Context) {
	formData := &model.FormData{}
	c.Bind(formData)

	if formData.Password == "" {
		c.HTML(http.StatusOK, "error_list", gin.H{"errors": &[]error{errors.New("Password is required")}})
		return
	}

	if errs := r.Service.Create(formData); errs != nil {
		c.HTML(http.StatusOK, "error_list", gin.H{"errors": errs})
	} else {
		c.Header("HX-Redirect", "/"+formData.Id)
	}
}
