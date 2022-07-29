package routes

import (
	controller "github.com/deep/controller"
	"github.com/deep/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/user/:user_id", controller.GetUser())
	incomingRoutes.PUT("/user/:user_id", controller.UpdateUser())
	incomingRoutes.DELETE("/user/:user_id", controller.DeleteUser())
}
