package usecase

import (
	"fmt"
	"net/http"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/helper"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/validator"
	v10 "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func usecaseValidation(ctx echo.Context, params any) *helper.ErrorStruct {
	if err := ctx.Validate(params); err != nil {
		errs := err.(v10.ValidationErrors)

		var newErr validator.ValidationError

		// errs[1]= v10.FieldError
		for _, val := range errs {
			message := fmt.Sprintf("validation error for '%s', Tag: %s", val.Field(), val.Tag())
			newErr.Message = append(newErr.Message, message)
		}

		return &helper.ErrorStruct{
			Code: http.StatusBadRequest,
			Err:  newErr,
		}
	}
	return nil
}
