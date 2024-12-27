package commands

import (
	"github.com/bbquite/go-pass-keeper/internal/cli/validator"
)

var (
	authParams = CommandParams{
		"username": {validateFunc: validator.StringValidation},
		"password": {validateFunc: validator.StringValidation},
	}

	pairParams = CommandParams{
		"key":  {validateFunc: validator.StringValidation},
		"pwd":  {validateFunc: validator.StringValidation},
		"meta": {validateFunc: validator.StringValidation},
	}

	textParams = CommandParams{
		"text": {validateFunc: validator.StringValidationUnlimit},
		"meta": {validateFunc: validator.StringValidation},
	}

	binaryParams = CommandParams{
		"filepath": {validateFunc: validator.FilePathValidation},
		"meta":     {validateFunc: validator.StringValidation},
	}

	cardParams = CommandParams{
		"number": {
			validateFunc: validator.CardNumberValidation,
			usage:        "ex. 4242 4242 4242 4242",
		},
		"cvv": {
			validateFunc: validator.CardCvvValidation,
			usage:        "ex. 777",
		},
		"owner": {
			validateFunc: validator.StringValidation,
			usage:        "ex. IVAN IVANOV",
		},
		"exp": {
			validateFunc: validator.DateValidation,
			usage:        "ex. 01.06",
		},
		"meta": {validateFunc: validator.StringValidation},
	}
)

func wrapIDParam(params CommandParams) CommandParams {
	params["id"] = CommandParamsItem{validateFunc: validator.IntValidation}
	return params
}
