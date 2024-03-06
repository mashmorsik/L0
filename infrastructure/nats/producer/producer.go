package producer

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

var SubjectNameCreated, _ = os.LookupEnv("JS_SUBJECT_NAME_CREATED")

type NatsProducer struct {
	streamContext nats.JetStreamContext
}

func NewNatsProducer(streamContext nats.JetStreamContext) *NatsProducer {
	return &NatsProducer{streamContext: streamContext}
}

// PublishOrders publishes messages to nats.JetStream.
func (n *NatsProducer) PublishOrders(fakeOrder string) {
	_, err := n.streamContext.Publish(SubjectNameCreated, []byte(fakeOrder))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Publisher  =>  Message:%s\n", fakeOrder)
	}
}
