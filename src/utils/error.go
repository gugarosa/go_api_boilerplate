package utils

import "log"

// LogError expects an dynamic number of error arguments,
// logs and returns the first error occurence
func LogError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

// LogFatalError expects an dynamic number of error arguments,
// logs and exits the system on first error occurence
func LogFatalError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
