package consumer

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/pkg/models"
	"github.com/nats-io/nats.go"
)

func ValidateMsg(msg nats.Msg) bool {
	order := models.Order{}
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Errf("can't unmarshal msg: %validate, err: %s", msg, err)
		return false
	}
	validate := validator.New(validator.WithRequiredStructEnabled())

	if validateOrderStruct(validate, order) != nil {
		log.Errf("invalid order struct: %v,err: %s", err)
		return false
	}

	return true
}

func validateOrderStruct(v *validator.Validate, order models.Order) error {
	err := v.Struct(order)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			log.Errf("invalid order struct: %v,err: %s", err)
			return err
		}
		return err
	}
	return nil
}
