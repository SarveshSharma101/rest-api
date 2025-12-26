package main

import (
	"rest-api/rest-api/internals/routes"
	"rest-api/rest-api/server"
)

func main() {
	server := server.GetServer()
	routes.UserRoutes(server)

	server.Run(":5000")
}
