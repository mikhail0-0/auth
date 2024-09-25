package router

import (
	"auth/src/config"
	"auth/src/controller"
	"auth/src/middleware"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.Default()

	router.POST(config.AUTH_PATH, controller.Authenticate)
	router.POST(config.AUTH_PATH+"/refresh", middleware.AuthGuard, controller.Refresh)

	router.GET("/protected", middleware.AuthGuard, controller.GetProtected)

	return router
}
