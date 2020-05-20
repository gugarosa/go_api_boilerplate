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

func getConfig() map[string]string {
	// Creating the configuration object
	config := map[string]string{
		"mode":      os.Getenv("MODE"),
		"dbUser":    os.Getenv("DB_USER"),
		"dbPass":    os.Getenv("DB_PASS"),
		"dbName":    os.Getenv("DB_NAME"),
		"dbPort":    os.Getenv("DB_PORT"),
		"redisPass": os.Getenv("REDIS_PASS"),
		"redisPort": os.Getenv("REDIS_PORT"),
	}

	return config
}

func main() {
	// Getting arguments from environment file
	c := getConfig()

	// Initializing the database and the cache
	db.InitDatabase(fmt.Sprintf("mongodb://%s:%s@db:%s", c["dbUser"], c["dbPass"], c["dbPort"]), c["dbName"])
	db.InitRedis(c["redisPort"], c["redisPass"])

	// Initializing the server
	server.InitServer(c["mode"])
}
