package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rest-api/rest-api/internals/routes"
	"rest-api/rest-api/server"
	"syscall"
	"time"
)

func main() {
	port := ":5000"
	router := server.GetServer()
	routes.UserRoutes(router)
	routes.AddGracefullRoute(router)
	routes.AddJobRoute(router)

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		fmt.Println("Server started on: ", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("err: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 11*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced shutdown")
	}

	log.Println("server exited gracefully")

}
