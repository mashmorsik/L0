package order

import (
	"encoding/json"
	"github.com/mashmorsik/L0/infrastructure/data"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/pkg/models"
	"github.com/nats-io/nats.go/jetstream"
)

type Order struct {
	Repo data.Data
}

func UnmarshalOrder(msg jetstream.Msg) (*models.Order, error) {
	order := models.Order{}
	err := json.Unmarshal(msg.Data(), &order)
	if err != nil {
		log.Errf("can't unmarshal msg: %v, err: %s", msg, err)
		return nil, err
	}

	return &order, nil
}

func (o *Order) AddNewOrder(msg jetstream.Msg) error {
	order, err := UnmarshalOrder(msg)
	if err != nil {
		log.Errf("can't unmarshal msg, err: %s", err)
	}

	err = o.Repo.AddOrder(*order)
	if err != nil {
		log.Errf("can't add order to database, err: %s", err)
	}

	return nil
}
