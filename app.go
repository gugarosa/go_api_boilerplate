package main

import (
	"vivere_api/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	config.ConfigureDatabase("mongodb://localhost:27017")

	// Initialize application
	r := gin.Default()

	// Adding router
	config.ConfigureRouter(r)

	r.Run()
}
