package handlers

import (
	"github.com/2hard4me/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/api")
	{
		api.GET("/get_batch/:id", h.GetBatch)
	}
	
	return router
}