package server

import (
	"github.com/gin-gonic/gin"
)

func GetServer() *gin.Engine {
	router := gin.Default()

	return router
}
