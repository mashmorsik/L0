package main

import (
	data2 "github.com/mashmorsik/L0/infrastructure/data"
	cache2 "github.com/mashmorsik/L0/infrastructure/data/cache"
	"github.com/mashmorsik/L0/infrastructure/nats"
	"github.com/mashmorsik/L0/infrastructure/nats/consumer"
	"github.com/mashmorsik/L0/infrastructure/nats/producer"
	"github.com/mashmorsik/L0/infrastructure/repository"
	"github.com/mashmorsik/L0/internal/order"
	log "github.com/mashmorsik/L0/pkg/logger"
)

func main() {
	log.BuildLogger()

	streamContext, err := nats.Connect()
	if err != nil {
		log.Errf("can't return stream context, err: %s", err)
		return
	}

	conn := data2.MustConnectPostgres()
	data2.MustMigrate(conn)

	orderRepo := repository.NewOrderRepo(data2.NewData(conn))
	createOrder := order.NewCreateOrder(orderRepo)

	producer.NewNatsProducer(streamContext).PublishOrders()
	consumer.NewNatsConsumer(createOrder)

	cache := cache2.NewCache()
	log.Infof("started cache: %v", cache)
}
