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

		if err != nil {
			log.Errf("consumer.Consume failed, err: %s", err)
		}

		if !ValidateMsg(*m) {
			return
		}

		md, err := m.Metadata()
		if err != nil {
			log.Errf("consumer.Consume failed, err: %s", err)
		}
		log.Infof("received msg on: subject:%s, md: %+v\n", m.Subject, md)

		//time.Sleep(time.Second * 1)

		o, err := order.UnmarshalOrder(*m)
		log.Infof("order: %v", o)

		err = n.o.CreateNewOrder(*o)
		if err != nil {
			log.Errf("can't add new o, err: %s", err)
			return
		}
	})

	if err != nil {
		log.Infof("Subscribe failed")
		return
	}

	//func Consumer(ctx context.Context, stream jetstream.Stream) error {
	//	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{Name: "Denis"})
	//	if err != nil {
	//		return err
	//	}
	//	log.Infof("consumer created: %+v\n", consumer)
	//
	//	cc, err := consumer.Consume(handleMsg)
	//	if err != nil {
	//		log.Errf("can't consume msg, err: %s", err)
	//	}
	//
	//	//add graceful shutdown
	//	defer cc.Stop()
	//
	//	return nil
	//}
	//
	//func handleMsg(msg jetstream.Msg) {
	//	err := msg.Ack()
	//	if err != nil {
	//		log.Errf("consumer.Consume failed, err: %s", err)
	//	}
	//
	//	if !ValidateMsg(msg) {
	//		return
	//	}
	//
	//	md, err := msg.Metadata()
	//	if err != nil {
	//		log.Errf("consumer.Consume failed, err: %s", err)
	//	}
	//	log.Infof("received msg on: subject:%s, md: %+v\n", msg.Subject(), md)
	//
	//	time.Sleep(time.Second * 1)
	//
	//	o, err := order.UnmarshalOrder(msg)
	//	data := data2.NewData(data2.MustConnectPostgres())
	//	err = data.AddOrderTx(*o)
	//	if err != nil {
	//		log.Errf("can't add new o, err: %s", err)
	//		return
	//	}
	//}
}
