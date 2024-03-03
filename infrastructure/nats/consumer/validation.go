package consumer

import (
	"encoding/json"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/pkg/models"
	"github.com/nats-io/nats.go"
)

func ValidateMsg(msg nats.Msg) bool {
	order := models.Order{}
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Errf("can't unmarshal msg: %v, err: %s", msg, err)
		return false
	}
	// sandbox to validate fields
	return true
}
