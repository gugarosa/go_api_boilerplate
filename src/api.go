package main

import (
	"fmt"
	"go_api_boilerplate/database"
	"go_api_boilerplate/server"
	"os"
)

func getConfig() map[string]string {
	// Creates the configuration object
	config := map[string]string{
		"mode":      os.Getenv("MODE"),
		"dbUser":    os.Getenv("DB_USER"),
		"dbPass":    os.Getenv("DB_PASS"),
		"dbName":    os.Getenv("DB_NAME"),
		"dbHost":    os.Getenv("DB_HOST"),
		"dbPort":    os.Getenv("DB_PORT"),
		"redisPass": os.Getenv("REDIS_PASS"),
		"redisHost": os.Getenv("REDIS_HOST"),
		"redisPort": os.Getenv("REDIS_PORT"),
	}

	return config
}

func main() {
	// Gets arguments from environment file
	c := getConfig()

	// Initializes the database and the cache
	database.InitMongo(fmt.Sprintf("mongodb://%s:%s@%s:%s", c["dbUser"], c["dbPass"], c["dbHost"], c["dbPort"]), c["dbName"])
	database.InitRedis(c["redisHost"], c["redisPort"], c["redisPass"])

	// Initializes the server
	server.InitServer(c["mode"])
}
