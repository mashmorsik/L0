package consumer

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/mashmorsik/L0/pkg/models"
	"github.com/nats-io/nats.go"
	errors2 "github.com/pkg/errors"
)

// ValidateMsg checks if the nats message is valid (can be unmarshaled into structure, all the required field are
// filled).
func ValidateMsg(msg nats.Msg) error {
	order := models.Order{}

	if err := json.Unmarshal(msg.Data, &order); err != nil {
		return errors2.WithMessagef(err, "can't unmarshal msg: %+v", msg)
	}

	if validate := validator.New(validator.WithRequiredStructEnabled()); validateOrderStruct(validate, order) != nil {
		return validateOrderStruct(validate, order)
	}

	return nil
}

// validateOrderStruct checks if all the required field in the models.Order structure are filled.
func validateOrderStruct(v *validator.Validate, order models.Order) error {
	if err := v.Struct(order); err != nil {
		var invalidValidationError *validator.InvalidValidationError
		var validationErrors *validator.ValidationErrors

		if errors.Is(err, invalidValidationError) {
			return errors2.WithMessage(err, "invalidValidationError")
		}
		if errors.Is(err, validationErrors) {
			return errors2.WithMessage(err, "validationErrors")
		}

		return err
	}
	return nil
}
