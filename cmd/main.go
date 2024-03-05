package main

import (
	"encoding/json"
	"fmt"
	"github.com/mashmorsik/L0/bench"
	data "github.com/mashmorsik/L0/infrastructure/data"
	"github.com/mashmorsik/L0/infrastructure/repository"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/pkg/models"
)

func main() {
	log.BuildLogger()

	//streamContext, err := nats.Connect()
	//if err != nil {
	//	log.Errf("can't return stream context, err: %s", err)
	//	return
	//}
	//fmt.Println(streamContext)
	//
	//js, err := nats.Conn()
	//if err != nil {
	//	log.Errf("can't return jet stream, err: %s", err)
	//	return
	//}
	//fmt.Println(js)
	//
	//ctx := context.Background()
	//orderCache := cache.NewOrderCache(ctx, time.Hour)

	conn := data.MustConnectPostgres()
	data.MustMigrate(conn)

	orderRepo := repository.NewOrderRepo(data.NewData(conn))

	generatedMsg := bench.GenerateMsg()
	msgMarsh, err := json.Marshal(generatedMsg)
	if err != nil {
		log.Errf("can't marshal json, err: %s", err)
	}

	fmt.Println(msgMarsh)

	var fakeOrder *models.Order

	err = json.Unmarshal(msgMarsh, fakeOrder)
	if err != nil {
		log.Errf("can't unmarshal msg, err: %s", err)
	}

	err = orderRepo.CreateOrder(*fakeOrder)
	if err != nil {
		log.Errf("can't add order to db, err: %s", err)
	}

	//createOrder := order.NewCreateOrder(orderRepo, orderCache)
	//err = createOrder.LoadCache()
	//if err != nil {
	//	log.Errf("can't load cache, err: %s", err)
	//	return
	//}
	//
	//go server.NewServer(createOrder).StartServer()
	//
	//consumer.NewNatsConsumer(createOrder).ConsumeOrders(ctx, streamContext)
}
