package main

import (
	"vivere_api/core"
	"vivere_api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initializing the database with desired arguments
	core.InitializeDatabase("mongodb://localhost:27017", "vivere")

	//
	middleware.InitializeRedis()

	// Initializing the application
	r := gin.Default()

	// Adding the router to the application
	core.AddRouter(r)

	// Running the application
	r.Run()
}
