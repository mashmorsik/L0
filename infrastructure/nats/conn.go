package nats

import (
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/nats-io/nats.go"
	jss "github.com/nats-io/nats.go/jetstream"
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

func Conn() (jss.JetStream, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	stream, err := jss.New(nc)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func JetStreamInit() (nats.JetStreamContext, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
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

	if stream == nil {
		log.Infof("Creating stream: %s\n", StreamName)

		stream, err = jetStream.AddStream(&nats.StreamConfig{
			Name:     StreamName,
			Subjects: []string{StreamSubjects},
		})

	}
	return nil
}
