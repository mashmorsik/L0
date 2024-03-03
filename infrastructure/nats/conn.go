package nats

import (
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/nats-io/nats.go"
)

const (
	StreamName     = "ORDER"
	StreamSubjects = "ORDER.*"
)

func Connect() (nats.JetStreamContext, error) {
	stream, err := JetStreamInit()
	if err != nil {
		log.Errf("can't init stream, err: %s", err)
	}
	err = CreateStream(stream)
	if err != nil {
		log.Errf("can't create stream, err: %s", err)
		return nil, err
	}
	return stream, nil
}

func JetStreamInit() (nats.JetStreamContext, error) {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	return js, nil
}

func CreateStream(jetStream nats.JetStreamContext) error {
	stream, err := jetStream.StreamInfo(StreamName)
	if err != nil {
		return err
	}

	// stream not found, create it
	if stream == nil {
		log.Infof("Creating stream: %s\n", StreamName)

		stream, err = jetStream.AddStream(&nats.StreamConfig{
			Name:     StreamName,
			Subjects: []string{StreamSubjects},
		})

	}
	return nil
}

//func Connect() error {
//	url := os.Getenv("NATS_URL")
//	if url == "" {
//		url = nats.DefaultURL
//	}
//
//	nc, err := nats.Connect(url)
//	if err != nil {
//		return err
//	}
//	defer func(nc *nats.Conn) {
//		if err = nc.Drain(); err != nil {
//			log.Errf("failed drain defer: %s", err)
//		}
//	}(nc)
//
//	js, err := jetstream.New(nc)
//	if err != nil {
//		return err
//	}
//
//	streamName := "maha"
//
//	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
//	defer cancel()
//
//	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
//		Name:     streamName,
//		Subjects: []string{"maha.>"},
//	})
//	if err != nil {
//		return err
//	}
//
//	log.Infof("created stream: %+v\n", stream)
//
//	err = producer.Producer(ctx, js)
//	if err != nil {
//		log.Errf("producer failed %s", err)
//		return err
//	}
//
//	err = consumer.Consumer(ctx, stream)
//	if err != nil {
//		log.Errf("consumer failed %s", err)
//		return err
//	}
//
//	return nil
//}
