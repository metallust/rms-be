package helper

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func ReadBody[T any](c *fiber.Ctx) (T, error) {
	var body T
	err := c.BodyParser(&body)
	if err != nil {
        log.Error("Error parsing body", err.Error())
		return body, errors.New(`invalid payload, please check your request and try again`)
	}

	err = ValidateStruct(body)
	if err != nil {
        log.Error("Error validating body", err.Error())
		return body, err
	}

	return body, nil
}

func ValidateStruct(s any) error {
	validate := validator.New()
	err := validate.Struct(s)
	errMsgs := make([]string, 0)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, err := range validationErrors {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"Error when validating %s: '%v'",
				err.Field(),
				err.Value(),
			))
		}
		return errors.New(errMsgs[0])
	}
	return nil
}
