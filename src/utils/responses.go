package utils

import "github.com/gin-gonic/gin"

// DynamicResponse expects a Gin context, a status identifier and a map
// to create the JSON response
func DynamicResponse(c *gin.Context, status int, m map[string]interface{}) {
	c.JSON(status, m)
}

// ConstantResponse expects a Gin context, a status identifier and a message
// to create the JSON response
func ConstantResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"response": message,
	})
}
