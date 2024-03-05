package producer

import (
	"github.com/nats-io/nats.go"
	"log"
)

const (
	SubjectNameReviewCreated = "ORDER.test"
)

type NatsProducer struct {
	streamContext nats.JetStreamContext
}

func NewNatsProducer(streamContext nats.JetStreamContext) *NatsProducer {
	return &NatsProducer{streamContext: streamContext}
}

func (n *NatsProducer) PublishOrders(fakeOrder string) {
	_, err := n.streamContext.Publish(SubjectNameReviewCreated, []byte(fakeOrder))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Publisher  =>  Message:%s\n", fakeOrder)
	}
}
