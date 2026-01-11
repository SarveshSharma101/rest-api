package routes

import (
	"rest-api/rest-api/internals/handler"

	"github.com/gin-gonic/gin"
)

func AddJobRoute(r *gin.Engine) {

	r.POST("/job", handler.Job)
	r.GET("/job/:jobId", handler.JobStatus)

}
