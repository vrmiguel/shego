package utils

import "log"

// AssertNonNil breaks execution if err is not nil
func AssertNonNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
