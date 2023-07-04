package utils

import (
	"context"

	"github.com/go-playground/validator/v10"
)

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate

type ErrorResponse struct {
	Message string                 `json:"message"`
	Errors  map[string]interface{} `json:"errors"`
}

func init() {
	validate = validator.New()
}

// Validate struct fields
func ValidateStruct(ctx context.Context, s interface{}) error {
	return validate.StructCtx(ctx, s)
}

func ShowErrors(err error) ErrorResponse {
	if err != nil {
		errorResponse := ErrorResponse{
			Message: "Validation error",
			Errors:  make(map[string]interface{}),
		}

		validationErrors := err.(validator.ValidationErrors)
		for _, validationErr := range validationErrors {
			field := validationErr.Field()
			actualTag := validationErr.ActualTag()
			errorResponse.Errors[field] = actualTag
		}
		return errorResponse
	}

	return ErrorResponse{
		Message: "No errors",
		Errors:  nil,
	}
}
