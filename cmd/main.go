package main

import (
	"context"
	"fmt"
	data "github.com/mashmorsik/L0/infrastructure/data"
	"github.com/mashmorsik/L0/infrastructure/data/cache"
	"github.com/mashmorsik/L0/infrastructure/nats"
	"github.com/mashmorsik/L0/infrastructure/nats/consumer"
	"github.com/mashmorsik/L0/infrastructure/repository"
	"github.com/mashmorsik/L0/infrastructure/server"
	"github.com/mashmorsik/L0/internal/order"
	log "github.com/mashmorsik/L0/pkg/logger"
	"time"
)

func main() {
	log.BuildLogger()

	streamContext, err := nats.Connect()
	if err != nil {
		log.Errf("can't return stream context, err: %s", err)
		return
	}
	fmt.Println(streamContext)

	ctx := context.Background()
	orderCache := cache.NewOrderCache(ctx, time.Hour)

	conn := data.MustConnectPostgres()
	data.MustMigrate(conn)

	orderRepo := repository.NewOrderRepo(data.NewData(conn))

	createOrder := order.NewCreateOrder(orderRepo, orderCache)
	err = createOrder.LoadCache()
	if err != nil {
		log.Errf("can't load cache, err: %s", err)
		return
	}

	go server.NewServer(createOrder).StartServer()

	consumer.NewNatsConsumer(createOrder).ConsumeOrders(ctx, streamContext)
}
