package consumer

import (
	"context"
	"fmt"
	"github.com/mashmorsik/L0/internal/order"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/nats-io/nats.go"
)

const (
	SubjectNameCreated = "WBORDER.test"
)

type NatsConsumer struct {
	o order.CreateOrder
}

func NewNatsConsumer(o order.CreateOrder) *NatsConsumer {
	return &NatsConsumer{o: o}
}

func (n *NatsConsumer) ConsumeOrders(ctx context.Context, js nats.JetStreamContext) {
	msgCh := make(chan *nats.Msg, 8192)
	_, err := js.ChanSubscribe(SubjectNameCreated, msgCh)
	if err != nil {
		log.Errf("Unable to ChanSubscribe, err: %s", err)
		return
	}
	for {
		select {
		case msg := <-msgCh:
			fmt.Println("[Received]", msg.Subject)
			n.orderHandler(msg)
		case <-ctx.Done():
			log.Errf("consumer ctx.Done, err: %s", err)
			return
		}
	}
}

func (n *NatsConsumer) orderHandler(m *nats.Msg) {
	err := m.Ack()
	if err != nil {
		log.Errf("Unable to Ack, err: %s", err)
		return
	}

	if !ValidateMsg(*m) {
		return
	}

	md, err := m.Metadata()
	if err != nil {
		log.Errf("Error extracting metadata, err: %s", err)
		return
	}
	log.Infof("received msg on: subject:%s, md: %+v\n", m.Subject, md)

	o, err := order.UnmarshalOrder(*m)
	if err != nil {
		log.Errf("Error unmarshaling order, err: %s", err)
		return
	}
	log.Infof("order: %v", o)

	err = n.o.CreateNewOrder(*o)
	if err != nil {
		log.Errf("can't add new o, err: %s", err)
		return
	}

	if err != nil {
		log.Errf("Subscribe failed, err: %s", err)
		return
	}
	return
}
