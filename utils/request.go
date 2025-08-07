package utils

import (
	"fmt"
	"incident-tracker/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindAndValidate[T any](c *gin.Context, req *T) errors.APIError {
	if err := c.ShouldBindJSON(req); err != nil {
		return errors.NewAPIError(400, "Invalid request body", err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMsg string
		for _, e := range validationErrors {
			errorMsg += fmt.Sprintf("Field '%s' failed on the '%s' tag. ", e.Field(), e.Tag())
		}
		return errors.NewAPIError(400, errorMsg, err)
	}

	return nil
}
