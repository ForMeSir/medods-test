package handler

import (
	"ruby/pkg/service"
	"ruby/pkg/user"

	"github.com/gin-gonic/gin"
)

//у этого сервера маршрутизатор будет от Джина(Gin)
		type Handler struct{
			services  *service.Service
			mongo     user.Storage
		}

		func NewHandler(services *service.Service, mongo user.Storage) *Handler{
			return &Handler{services: services, mongo:mongo}
		}

		func (h *Handler) InitRoutes() *gin.Engine{
		router := gin.New()
		auth := router.Group("/auth")
			{
				auth.POST("/sign-in", h.signIn)
				auth.POST("/refresh", h.refresh)
			} 
		return router
		}