package nats

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/nats-io/nats.go"
	"os"
)

var (
	StreamName, _     = os.LookupEnv("JS_STREAM_NAME")
	StreamSubjects, _ = os.LookupEnv("JS_STREAM_SUBJECTS")
)

// Connect initiates the connection to the nats system and creates new stream if the stream with a defined name is not
// found.
func Connect(ctx context.Context) (nats.JetStreamContext, error) {
	stream, err := JetStreamInit(ctx)
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

// JetStreamInit initiates the connection to the nats system.
func JetStreamInit(ctx context.Context) (nats.JetStreamContext, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		nc.Close()
	}()

	return js, nil
}

// CreateStream creates new stream if the stream with a defined name is not found.
func CreateStream(jetStream nats.JetStreamContext) error {
	stream, err := jetStream.StreamInfo(StreamName)
	if err != nil {
		log.Errf("can't get stream info, err: %s", err)
		return err
	}

	if stream == nil {
		log.Infof("Creating stream: %s\n", StreamName)

		stream, err = jetStream.AddStream(&nats.StreamConfig{
			Name:     StreamName,
			Subjects: []string{StreamSubjects},
		})
		if err != nil {
			log.Errf("can't create stream, err: %s", err)
			return err
		}
	}
	return nil
}
