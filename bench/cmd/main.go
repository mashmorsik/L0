package main

import (
	"encoding/json"
	"fmt"
	"github.com/mashmorsik/L0/bench"
	"github.com/mashmorsik/L0/infrastructure/nats"
	"github.com/mashmorsik/L0/infrastructure/nats/producer"
	log "github.com/mashmorsik/L0/pkg/logger"
)

func main() {
	newFakeMsg := bench.GenerateMsg()

	msgJSON, err := json.Marshal(newFakeMsg)
	if err != nil {
		log.Errf("can't marshal fake msg into JSON, err: %s", err)
		return
	}

	streamContext, err := nats.Connect()
	if err != nil {
		log.Errf("can't return stream context, err: %s", err)
		return
	}
	fmt.Println(streamContext)

	producer.NewNatsProducer(streamContext).PublishOrders(string(msgJSON))
}
