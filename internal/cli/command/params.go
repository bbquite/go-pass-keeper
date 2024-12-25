package command

import "github.com/bbquite/go-pass-keeper/internal/cli/validator"

var authParams = CommandParams{
	"username": {validateFunc: validator.StringValidation},
	"password": {validateFunc: validator.StringValidation},
}

var onlyIDParams = CommandParams{
	"id": {validateFunc: validator.IntValidation},
}

var pairParams = CommandParams{
	"key":  {validateFunc: validator.StringValidation},
	"pwd":  {validateFunc: validator.StringValidation},
	"meta": {validateFunc: validator.StringValidation},
}

var textParams = CommandParams{
	"text": {validateFunc: validator.StringValidationUnlimit},
	"meta": {validateFunc: validator.StringValidation},
}

var binaryParams = CommandParams{
	"filepath": {validateFunc: validator.FilePathValidation},
	"meta":     {validateFunc: validator.StringValidation},
}

var cardParams = CommandParams{
	"number": {
		validateFunc: validator.CardNumberValidation,
		usage:        "4242 4242 4242 4242",
	},
	"cvv": {
		validateFunc: validator.CardCvvValidation,
		usage:        "777",
	},
	"owner": {
		validateFunc: validator.StringValidation,
		usage:        "IVAN IVANOV",
	},
	"exp": {
		validateFunc: validator.DateValidation,
		usage:        "ex. 01.06",
	},
	"meta": {validateFunc: validator.StringValidation},
}

func wrapIDParam(params CommandParams) CommandParams {
	params["id"] = CommandParamsItem{validateFunc: validator.IntValidation}
	return params
}
