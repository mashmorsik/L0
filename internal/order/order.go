package order

import (
	"context"
	"encoding/json"
	"github.com/mashmorsik/L0/infrastructure/data/cache"
	"github.com/mashmorsik/L0/infrastructure/repository"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/pkg/models"
	"github.com/nats-io/nats.go"
	"time"
)

type CreateOrder struct {
	Ctx   context.Context
	Repo  repository.Repository
	Cache *cache.OrderCache
}

func NewCreateOrder(ctx context.Context, repo repository.Repository, orderCache cache.OrderCache) CreateOrder {
	return CreateOrder{Ctx: ctx, Repo: repo, Cache: &orderCache}
}

// UnmarshalOrder umnmarshals nats.Msg into an order struct.
func UnmarshalOrder(msg nats.Msg) (*models.Order, error) {
	order := models.Order{}
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Errf("can't unmarshal msg: %v, err: %s", msg, err)
		return nil, err
	}

	return &order, nil
}

// CreateNewOrder adds a new order to the database and local cache.
func (c *CreateOrder) CreateNewOrder(order models.Order) error {
	_, cancel := context.WithTimeout(c.Ctx, time.Second*5)
	defer cancel()

	err := c.Repo.CreateOrder(order)
	if err != nil {
		log.Errf("can't create order, err: %s", err)
		return nil
	}

	c.Cache.Set(order.OrderUid, order, time.Hour)

	return nil
}

// CacheWarmUp fills local cache with orders from the database when the app starts.
func (c *CreateOrder) CacheWarmUp() error {
	_, cancel := context.WithTimeout(c.Ctx, time.Second*5)
	defer cancel()

	orders, err := c.Repo.GetOrders()
	if err != nil {
		log.Errf("can't get orders from db, err: %s", err)
		return err
	}

	for _, o := range orders {
		c.Cache.Set(o.OrderUid, *o, time.Hour)
		log.Infof("added to cache order_id: %s", o.OrderUid)
	}

	return nil
}

// GetOrderFromCache returns the order from the local cache by its id.
func (c *CreateOrder) GetOrderFromCache(orderID string) (*models.Order, error) {
	order, ok := c.Cache.Get(orderID)
	if !ok {
		log.Infof("no order id in cache, orderID: %s", orderID)
		return nil, nil
	}

	return order, nil
}
