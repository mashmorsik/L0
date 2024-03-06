package main

import (
	"context"
	data "github.com/mashmorsik/L0/infrastructure/data"
	"github.com/mashmorsik/L0/infrastructure/data/cache"
	"github.com/mashmorsik/L0/infrastructure/nats"
	"github.com/mashmorsik/L0/infrastructure/nats/consumer"
	"github.com/mashmorsik/L0/infrastructure/repository"
	"github.com/mashmorsik/L0/infrastructure/server"
	"github.com/mashmorsik/L0/internal/order"
	log "github.com/mashmorsik/L0/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.BuildLogger()

	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGKILL)

	go func() {
		<-sigCh
		log.Infof("context done")
		cancel()
	}()

	streamContext, err := nats.Connect(ctx)
	if err != nil {
		log.Errf("can't return stream context, err: %s", err)
		return
	}
	log.Infof("connected stream context: %+v", streamContext)

	orderCache := cache.NewOrderCache(ctx, time.Second)

	conn := data.MustConnectPostgres(ctx)
	data.MustMigrate(conn)

	orderRepo := repository.NewOrderRepo(ctx, data.NewData(ctx, conn))
	createOrder := order.NewCreateOrder(ctx, orderRepo, orderCache)

	err = createOrder.CacheWarmUp()
	if err != nil {
		log.Errf("can't load cache, err: %s", err)
		return
	}

	go func() {
		httpServer := server.NewServer(createOrder)
		if err = httpServer.StartServer(); err != nil {
			log.Err(err, "start http server failed")
		}
	}()

	consumer.NewNatsConsumer(ctx, createOrder).ConsumeOrders(streamContext)
}
