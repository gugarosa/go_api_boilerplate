package main

import (
	"vivere_api/db"
	"vivere_api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initializing the database with desired arguments
	db.InitializeDatabase("mongodb://localhost:27017", "vivere")

	// Initializing the application
	r := gin.Default()

	// Adding the router to the application
	handlers.AddRouter(r)

	// Running the application
	r.Run()
}
