package main

import (
	"fmt"
	"os"
	"vivere_api/db"
	"vivere_api/server"
	"vivere_api/utils"

	"github.com/joho/godotenv"
)

func init() {
	// Loading environment file
	err := godotenv.Load()

	// Handles a possible fatal error
	utils.HandleFatalError(err)
}

func main() {
	// Getting arguments from environment file
	mode := os.Getenv("MODE")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Initializing the database
	db.InitDatabase(fmt.Sprintf("mongodb://%s:%s@db:%s", dbUser, dbPass, dbPort), dbName)

	// Initializing middlewares
	// middleware.InitRedis()

	// Initializing the server
	server.InitServer(mode)
}
