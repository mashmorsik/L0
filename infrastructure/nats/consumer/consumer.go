package consumer

import (
	"context"
	"fmt"
	data2 "github.com/mashmorsik/L0/infrastructure/data"
	"github.com/mashmorsik/L0/internal/order"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

func Consumer(ctx context.Context, stream jetstream.Stream) error {
	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{Name: "Denis"})
	if err != nil {
		return err
	}
	fmt.Printf("consumer created: %+v\n", consumer)

	cc, err := consumer.Consume(handleMsg)

	//add graceful shutdown
	defer cc.Stop()

	return nil
}

func handleMsg(msg jetstream.Msg) {
	err := msg.Ack()
	if err != nil {
		fmt.Println("consumer.Consume failed", err)
	}

	if !ValidateMsg(msg) {
		return
	}

	md, err := msg.Metadata()
	if err != nil {
		fmt.Println("consumer.Consume failed", err)
	}
	fmt.Printf("received msg on: subject:%s, md: %+v\n", msg.Subject(), md)

	time.Sleep(time.Second * 1)

	o, err := order.UnmarshalOrder(msg)
	data := data2.NewData(data2.MustConnectPostgres())
	err = data.AddOrder(*o)
	if err != nil {
		log.Errf("can't add new o, err: %s", err)
		return
	}
}
