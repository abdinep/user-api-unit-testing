package router

import (
	"project/controller"

	"github.com/gin-gonic/gin"
)

func Ginsetup() *gin.Engine {
	server := gin.Default()
	UserRoute(server)
	return server
}
func UserRoute(router *gin.Engine) {
	router.POST("/signup", controller.Signup)
}
