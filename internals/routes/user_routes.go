package routes

import (
	"rest-api/rest-api/internals/handler"
	middleware "rest-api/rest-api/internals/middleware/auth"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {

	r.POST("/user", handler.CreateUser)
	r.POST("/login", handler.Login)
	r.GET("/refresh", handler.Refresh)

	nr := r.Group("/")
	nr.Use(middleware.JwtAuth())
	{
		nr.GET("/user", handler.GetAllUser)
		nr.GET("/user/:id", handler.GetUserById)
		nr.DELETE("/user/:id", handler.Delete)
		nr.PUT("/user/:id", handler.Update)
		nr.PATCH("/user/:id", handler.Patch)
	}
}
