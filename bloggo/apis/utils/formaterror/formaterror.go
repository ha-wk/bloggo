package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "nickname") {
		return errors.New("NICKNAME ALREADY TAKEN")
	}
	if strings.Contains(err, "EMAIL") {
		return errors.New("EMAIL ALREADY TAKEN")
	}
	if strings.Contains(err, "TITLE") {
		return errors.New("TITLE ALREADY TAKEN")
	}
	if strings.Contains(err, "HASHED PASSWORD") {
		return errors.New("INCORRECT PASSWORD")

	}
	return errors.New("INCORRECT DETAILS")
}

//when we wire the controllers and routes package this will make sense
