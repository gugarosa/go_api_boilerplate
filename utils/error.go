package utils

import "log"

// HandleError expects an dynamic number of error arguments,
// logs if an error has happened and returns `false`, otherwise returns `true`
func HandleError(errs ...error) bool {
	// For every possible error
	for _, err := range errs {
		// Checks if it exists, logs it and return as `false`
		if err != nil {
			log.Println(err)

			return false
		}
	}

	return true
}

// HandleFatalError expects an dynamic number of error arguments,
// logs and exits if a fatal error has happened
func HandleFatalError(errs ...error) bool {
	// For every possible error
	for _, err := range errs {
		// Checks if it exists, logs it and exits the system
		if err != nil {
			log.Fatal(err)
		}
	}

	return true
}
