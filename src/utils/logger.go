package utils

import "log"

// LogError expects a dynamic number of error arguments,
// logs and returns the first error occurrence
func LogError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

// LogFatalError expects a dynamic number of error arguments,
// logs and exits the system on first error occurrence
func LogFatalError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
