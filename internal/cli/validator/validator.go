package validator

import "fmt"

type ValidateFunc func(string) error

func BaseStringValidation(param string) error {
	if param == "" || len(param) > 250 {
		return fmt.Errorf("invalid string length")
	}
	return nil
}

func FilePathValidation(param string) error {
	if param == "" || len(param) > 250 {
		return fmt.Errorf("invalid string length")
	}
	return nil
}
