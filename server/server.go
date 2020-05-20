package server

import "github.com/gin-gonic/gin"

// InitServer expects an initialization mode to start the application
func InitServer(mode string) {
	// Setting application mode and creating it
	gin.SetMode(mode)
	r := gin.Default()

	// Initializing the router
	InitRouter(r)

	// Running the application
	r.Run()
}
