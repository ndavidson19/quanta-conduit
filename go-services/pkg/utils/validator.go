package utils

import (
    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func GetValidator() *validator.Validate {
    if validate == nil {
        validate = validator.New()
    }
    return validate
}

func ValidateStruct(s interface{}) error {
    return GetValidator().Struct(s)
}