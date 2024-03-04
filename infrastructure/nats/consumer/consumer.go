package consumer

import (
	"github.com/mashmorsik/L0/internal/order"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/nats-io/nats.go"
)

const (
	SubjectNameReviewCreated = "ORDER.test"
)

type NatsConsumer struct {
	o order.CreateOrder
}

func NewNatsConsumer(o order.CreateOrder) *NatsConsumer {
	return &NatsConsumer{o: o}
}

func (n *NatsConsumer) ConsumeOrders(js nats.JetStreamContext) {
	_, err := js.Subscribe(SubjectNameReviewCreated, func(m *nats.Msg) {
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

		//time.Sleep(time.Second * 1)

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
	})

	if err != nil {
		log.Errf("Subscribe failed, err: %s", err)
		return
	}
	return
}
