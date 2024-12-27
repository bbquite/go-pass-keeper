package validator

import (
	"os"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	validationIS "github.com/go-ozzo/ozzo-validation/is"
)

type ValidateFunc func(string) error

func IntValidation(param string) error {
	return validation.Validate(
		param,
		validation.Required,
		validationIS.Digit,
	)
}

func StringValidation(param string) error {
	return validation.Validate(
		param,
		validation.Required,
		validation.Length(1, 250),
	)
}

func StringValidationUnlimit(param string) error {
	return validation.Validate(
		param,
		validation.Required,
	)
}

func FilePathValidation(param string) error {
	_, err := os.ReadFile(param)
	if err != nil {
		return err
	}
	return nil
}

func CardNumberValidation(param string) error {
	return validation.Validate(
		param,
		validation.Required,
		validationIS.CreditCard,
	)
}

func CardCvvValidation(param string) error {
	return validation.Validate(
		param,
		validation.Required,
		validation.Match(regexp.MustCompile("^[0-9]{3}$")),
	)
}

func DateValidation(param string) error {
	return validation.Validate(
		param,
		validation.Required,
		validation.Date("01.06"),
	)
}
