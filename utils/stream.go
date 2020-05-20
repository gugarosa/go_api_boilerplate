package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

// HandleError expects an dynamic number of error arguments,
// logs and returns the first error occurence
func HandleError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

// HandleFatalError expects an dynamic number of error arguments,
// logs and exits the system on first error occurence
func HandleFatalError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

// SendStaticResponse expects a Gin context, status identifier and message
// to create the JSON response
func SendStaticResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status":  status,
		"message": message,
	})
}
