package consumer

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mashmorsik/L0/internal/order"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"os"
)

var SubjectNameCreated, _ = os.LookupEnv("JS_SUBJECT_NAME_CREATED")

type NatsConsumer struct {
	Ctx context.Context
	o   order.CreateOrder
}

func NewNatsConsumer(ctx context.Context, o order.CreateOrder) *NatsConsumer {
	return &NatsConsumer{Ctx: ctx, o: o}
}

// ConsumeOrders subscribes to the nats.JetStream, reads and handles messages.
func (n *NatsConsumer) ConsumeOrders(js nats.JetStreamContext) error {
	msgCh := make(chan *nats.Msg, 100)
	_, err := js.ChanSubscribe(SubjectNameCreated, msgCh)
	if err != nil {
		return errors.WithMessagef(err, "ChanSubscribe failed")
	}
	for {
		select {
		case msg := <-msgCh:
			log.Infof("[Received] subject: %s", msg.Subject)

			if err = n.orderHandler(msg); err != nil {
				log.Err(err, "handle stream message failed")
			}

		case <-n.Ctx.Done():
			log.Warn("consumer Ctx.Done")
			return nil
		}
	}
}

// orderHandler unmarshals nats messages and call CreateNewOrder function that adds the order to the database and local
// cache.
func (n *NatsConsumer) orderHandler(m *nats.Msg) error {
	err := m.Ack()
	if err != nil {
		return errors.WithMessagef(err, "unable to Ack nats message")
	}

	if ValidateMsg(*m) != nil {
		return errors.WithMessagef(err, "unable to Ack nats message")
	}

	md, err := m.Metadata()
	if err != nil {
		log.Errf("Error extracting metadata, err: %s", err)
		return err
	}
	log.Infof("received msg on: subject:%s, md: %+v\n", m.Subject, md)

	o, err := order.UnmarshalOrder(*m)
	if err != nil {
		log.Errf("Error unmarshaling order, err: %s", err)
		return err
	}
	log.Infof("order: %v", o)

	err = n.o.CreateNewOrder(*o)
	if err != nil {
		log.Errf("can't add new o, err: %s", err)
		return err
	}

	if err != nil {
		log.Errf("Subscribe failed, err: %s", err)
		return err
	}

	return nil
}
