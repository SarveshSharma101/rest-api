package routes

import (
	"rest-api/rest-api/internals/handler"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/user", handler.GetAllUser)
	r.GET("/user/:id", handler.GetUserById)
	r.POST("/user", handler.CreateUser)
	r.DELETE("/user/:id", handler.Delete)
	r.PUT("/user/:id", handler.Update)
	r.PATCH("/user/:id", handler.Patch)
}
