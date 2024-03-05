package producer

import (
	"github.com/nats-io/nats.go"
	"log"
)

const (
	SubjectNameCreated = "WBORDER.test"
)

type NatsProducer struct {
	streamContext nats.JetStreamContext
}

func NewNatsProducer(streamContext nats.JetStreamContext) *NatsProducer {
	return &NatsProducer{streamContext: streamContext}
}

func (n *NatsProducer) PublishOrders(fakeOrder string) {
	_, err := n.streamContext.Publish(SubjectNameCreated, []byte(fakeOrder))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Publisher  =>  Message:%s\n", fakeOrder)
	}
}
