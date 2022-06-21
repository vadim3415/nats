package handler

import (
	"github.com/gin-gonic/gin"
	"nats/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	AllGroup := router.Group("/")
	{
		AllGroup.GET("/order/:id", h.getId)
		AllGroup.GET("/sub", h.natsSub)
		AllGroup.GET("/pub", h.natsPub)
	}
	return router
}
