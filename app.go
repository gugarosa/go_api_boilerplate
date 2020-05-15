package main

import (
	"fmt"
	"log"
	"os"
	"vivere_api/core"
	"vivere_api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Loading environment file
	err := godotenv.Load()

	// If environment file could not be loaded
	if err != nil {
		// Invokes a fatal error
		log.Fatal("The .env file could not be found.")
	}
}

func main() {
	// Getting arguments from environment file
	mode := os.Getenv("MODE")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Initializing the database with desired arguments
	core.InitializeDatabase(fmt.Sprintf("mongodb://%s:%s@db:%s", dbUser, dbPass, dbPort), dbName)

	//
	middleware.InitializeRedis()

	// Setting application mode and initializing it
	gin.SetMode(mode)
	r := gin.Default()

	// Adding the router to the application
	core.AddRouter(r)

	// Running the application
	r.Run()
}
