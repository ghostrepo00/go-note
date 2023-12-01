package app

import (
	"net/http"

	"github.com/ghostrepo00/go-note/config"
	"github.com/gin-gonic/gin"
)

type webHandler struct {
	AppConfig *config.AppConfig
	Service   *appService
}

func NewWebHandler(appConfig *config.AppConfig, service *appService) *webHandler {
	return &webHandler{appConfig, service}
}

func (r *webHandler) Default(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}
