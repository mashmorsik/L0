package nats

import (
	"context"
	"fmt"
	"github.com/mashmorsik/L0/infrastructure/nats/consumer"
	"github.com/mashmorsik/L0/infrastructure/nats/producer"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"os"
	"time"
)

func Connect() error {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return err
	}
	defer func(nc *nats.Conn) {
		if err = nc.Drain(); err != nil {
			log.Errf("failed drain defer: %s", err)
		}
	}(nc)

	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}

	streamName := "maha"

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"maha.>"},
	})
	if err != nil {
		return err
	}

	fmt.Printf("created stream: %+v\n", stream)

	err = producer.Producer(ctx, js)
	if err != nil {
		log.Errf("producer failed %s", err)
		return err
	}

	err = consumer.Consumer(ctx, stream)
	if err != nil {
		log.Errf("consumer failed %s", err)
		return err
	}

	return nil
}
