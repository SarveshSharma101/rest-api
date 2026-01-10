package routes

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func AddGracefullRoute(r *gin.Engine) {
	r.GET("/handleme", func(ctx *gin.Context) {
		rc := ctx.Request.Context()
		select {
		case <-rc.Done():
			fmt.Println("Server shutdown")
		default:
			fmt.Println("sleeping")
			time.Sleep(10 * time.Second)
			fmt.Println("Done")
		}
	})
}
