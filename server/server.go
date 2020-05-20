package server

import "github.com/gin-gonic/gin"

// InitServer expects an initialization mode to start the application
func InitServer(mode string) {
	// Sets application mode and creates it
	gin.SetMode(mode)
	r := gin.Default()

	// Initializes the router
	InitRouter(r)

	// Runs the application
	r.Run()
}
