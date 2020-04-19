package utils

import "log"

// HandleError expects an dynamic number of error arguments,
// logs if an error has happened and returns `true`
func HandleError(errs ...error) bool {
	// For every possible error
	for _, err := range errs {
		// Check if it exists, logs it and return as `true`
		if err != nil {
			log.Println(err)
			return true
		}
	}
	return false
}
