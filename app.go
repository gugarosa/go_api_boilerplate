package main

import (
	"vivere_api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/get", handlers.LoginHandle)

	r.Run()
}
