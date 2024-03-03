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

func (n *NatsProducer) PublishOrders() {
	jsonData := `
		{
		"order_uid": "ff563feb7b2b84b6test",
		"track_number": "WBILMTESTTRACK",
		"entry": "WBIL",
		"delivery": {
		 "name": "Test Testof",
		 "phone": "+9720000000",
		 "zip": "2639809",
		 "city": "Kiryat Mozkin",
		 "address": "Ploshad Mira 15",
		 "region": "Kraiot",
		 "email": "test@gmail.com"
		},
		"payment": {
		 "transaction": "ff563feb7b2b84b6test",
		 "request_id": "",
		 "currency": "USD",
		 "provider": "wbpay",
		 "amount": 1817,
		 "payment_dt": 1637907727,
		 "bank": "alpha",
		 "delivery_cost": 1500,
		 "goods_total": 317,
		 "custom_fee": 0
		},
		"items": [
		 {
		   "chrt_id": 9934930,
		   "track_number": "WBILMTESTTRACK",
		   "price": 453,
		   "rid": "ab4219087a764ae0btest",
		   "name": "Mascaras",
		   "sale": 30,
		   "size": "0",
           "count": 1,
		   "total_price": 317,
		   "nm_id": 2389212,
		   "brand": "Vivienne Sabo",
		   "status": 202
		 }
		],
		"locale": "en",
		"internal_signature": "",
		"customer_id": "test",
		"delivery_service": "meest",
		"shardkey": "9",
		"sm_id": 99,
		"date_created": "2021-11-26T06:22:19Z",
		"oof_shard": "1"
		}`
	_, err := n.streamContext.Publish(SubjectNameReviewCreated, []byte(jsonData))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Publisher  =>  Message:%s\n", jsonData)
	}
}

//func Producer(ctx context.Context, js jetstream.JetStream) error {
//	jsonData := `
//	{
//	"order_uid": "b563feb7b2b84b6test",
//	"track_number": "WBILMTESTTRACK",
//	"entry": "WBIL",
//	"delivery": {
//	 "name": "Test Testov",
//	 "phone": "+9720000000",
//	 "zip": "2639809",
//	 "city": "Kiryat Mozkin",
//	 "address": "Ploshad Mira 15",
//	 "region": "Kraiot",
//	 "email": "test@gmail.com"
//	},
//	"payment": {
//	 "transaction": "b563feb7b2b84b6test",
//	 "request_id": "",
//	 "currency": "USD",
//	 "provider": "wbpay",
//	 "amount": 1817,
//	 "payment_dt": 1637907727,
//	 "bank": "alpha",
//	 "delivery_cost": 1500,
//	 "goods_total": 317,
//	 "custom_fee": 0
//	},
//	"items": [
//	 {
//	   "chrt_id": 9934930,
//	   "track_number": "WBILMTESTTRACK",
//	   "price": 453,
//	   "rid": "ab4219087a764ae0btest",
//	   "name": "Mascaras",
//	   "sale": 30,
//	   "size": "0",
//	   "total_price": 317,
//	   "nm_id": 2389212,
//	   "brand": "Vivienne Sabo",
//	   "status": 202
//	 }
//	],
//	"locale": "en",
//	"internal_signature": "",
//	"customer_id": "test",
//	"delivery_service": "meest",
//	"shardkey": "9",
//	"sm_id": 99,
//	"date_created": "2021-11-26T06:22:19Z",
//	"oof_shard": "1"
//	}`
//
//	publish1, err := js.Publish(ctx, streamSubject, []byte(jsonData))
//
//	if err != nil {
//		return err
//	}
//	log.Infof("publish: %+v\n", publish1)
//	return nil
//}
