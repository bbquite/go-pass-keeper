package validator

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	validationIS "github.com/go-ozzo/ozzo-validation/is"
)

type ValidateFunc func(string) error

func IntValidation(param string) error {
	err := validation.Validate(
		param,
		validation.Required,
		validationIS.Digit,
	)
	if err != nil {
		return err
	}

	return nil
}

func StringValidation(param string) error {
	err := validation.Validate(
		param,
		validation.Required,
		validation.Length(1, 250),
	)
	if err != nil {
		return err
	}

	return nil
}

func StringValidationUnlimit(param string) error {
	err := validation.Validate(
		param,
		validation.Required,
	)
	if err != nil {
		return err
	}

	return nil
}

func FilePathValidation(param string) error {
	if param == "" || len(param) > 250 {
		return fmt.Errorf("invalid string length")
	}
	return nil
}

func CardNumberValidation(param string) error {
	err := validation.Validate(
		param,
		validation.Required,
		validationIS.CreditCard,
	)
	if err != nil {
		return err
	}

	return nil

}

func CardCvvValidation(param string) error {
	err := validation.Validate(
		param,
		validation.Required,
		validation.Match(regexp.MustCompile("^[0-9]{3}$")),
	)
	if err != nil {
		return err
	}
	return nil
}

func DateValidation(param string) error {
	err := validation.Validate(
		param,
		validation.Required,
		validation.Date("01.06"),
	)
	if err != nil {
		return err
	}

	return nil
}
