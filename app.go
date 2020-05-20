package main

import (
	"os"
	"vivere_api/db"
	"vivere_api/server"
	"vivere_api/utils"

	"github.com/joho/godotenv"
)

func init() {
	// Loads environment file
	loadErr := godotenv.Load()

	// Handles a possible fatal error
	utils.HandleFatalError(loadErr)
}

func getConfig() map[string]string {
	// Creates the configuration object
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
	// Gets arguments from environment file
	c := getConfig()

	// Initializes the database and the cache
	// db.InitDatabase(fmt.Sprintf("mongodb://%s:%s@db:%s", c["dbUser"], c["dbPass"], c["dbPort"]), c["dbName"])
	db.InitDatabase("mongodb://localhost:27017", c["dbName"])
	// db.InitRedis(c["redisPort"], c["redisPass"])

	// Initializes the server
	server.InitServer(c["mode"])
}
