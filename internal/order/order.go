package order

import (
	"encoding/json"
	"github.com/mashmorsik/L0/infrastructure/repository"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/pkg/models"
	"github.com/nats-io/nats.go"
)

type CreateOrder struct {
	Repo repository.Repository
}

func NewCreateOrder(repo repository.Repository) CreateOrder {
	return CreateOrder{Repo: repo}
}

func UnmarshalOrder(msg nats.Msg) (*models.Order, error) {
	order := models.Order{}
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Errf("can't unmarshal msg: %v, err: %s", msg, err)
		return nil, err
	}

	return &order, nil
}

func (c *CreateOrder) CreateNewOrder(order models.Order) error {
	err := c.Repo.CreateOrder(order)
	if err != nil {
		log.Errf("can't create order, err: %s", err)
		return nil
	}

	return nil
}
