package order

import (
	"encoding/json"
	"fmt"
	"github.com/mashmorsik/L0/infrastructure/data/cache"
	"github.com/mashmorsik/L0/infrastructure/repository"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/pkg/models"
	"github.com/nats-io/nats.go"
	"time"
)

type CreateOrder struct {
	Repo  repository.Repository
	Cache *cache.OrderCache
}

func NewCreateOrder(repo repository.Repository, orderCache cache.OrderCache) CreateOrder {
	return CreateOrder{Repo: repo, Cache: &orderCache}
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
		//нужно ли возвращать здесь ошибку?
		return nil
	}

	c.Cache.Set(order.OrderUid, order, time.Hour)
	log.Infof("added order: %s to cache", order.OrderUid)
	fmt.Println(c.Cache.Get(order.OrderUid))

	time.Sleep(time.Second * 5)

	return nil
}

func (c *CreateOrder) LoadCache() error {
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

func (c *CreateOrder) GetOrderFromCache(orderID string) (*models.Order, error) {
	order, ok := c.Cache.Get(orderID)
	if !ok {
		log.Infof("no order in cache, orderID: %s", orderID)
		return nil, nil
	}

	return order, nil
}
