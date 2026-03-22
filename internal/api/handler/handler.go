package handler

import (
	"WBTech_L3.2/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type Handler struct {
	services *service.Service
	logger   zlog.Zerolog
}

func NewHandler(services *service.Service, logger zlog.Zerolog) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *ginext.Engine {
	r := ginext.New("")

	r.POST("/shorten", handlerFunc(h.handleCreate))
	r.GET("/s/:short_url", handlerFunc(h.handleRedirect))
	r.GET("/analytics/:short_url", handlerFunc(h.handleGetStats))

	r.Static("/static", "./web")
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	return r
}
