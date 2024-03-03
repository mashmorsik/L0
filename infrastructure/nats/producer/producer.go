package producer

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/nats-io/nats.go/jetstream"
)

type Test struct {
	UserId string `json:"user_id"`
}

func Producer(ctx context.Context, js jetstream.JetStream) error {
	newId := Test{UserId: "379548347"}
	jsonData, err := json.Marshal(newId)
	if err != nil {
		log.Errf("can't marshal struct: %s", newId)
	}

	fmt.Println(jsonData)

	publish1, err := js.Publish(ctx, "maha.test", []byte("New msg"))
	//publish1, err := js.Publish(ctx, "maha.test", []byte(`
	//{
	//"order_uid": "b563feb7b2b84b6test",
	//"track_number": "WBILMTESTTRACK",
	//"entry": "WBIL",
	//"delivery": {
	//  "name": "Test Testov",
	//  "phone": "+9720000000",
	//  "zip": "2639809",
	//  "city": "Kiryat Mozkin",
	//  "address": "Ploshad Mira 15",
	//  "region": "Kraiot",
	//  "email": "test@gmail.com"
	//},
	//"payment": {
	//  "transaction": "b563feb7b2b84b6test",
	//  "request_id": "",
	//  "currency": "USD",
	//  "provider": "wbpay",
	//  "amount": 1817,
	//  "payment_dt": 1637907727,
	//  "bank": "alpha",
	//  "delivery_cost": 1500,
	//  "goods_total": 317,
	//  "custom_fee": 0
	//},
	//"items": [
	//  {
	//    "chrt_id": 9934930,
	//    "track_number": "WBILMTESTTRACK",
	//    "price": 453,
	//    "rid": "ab4219087a764ae0btest",
	//    "name": "Mascaras",
	//    "sale": 30,
	//    "size": "0",
	//    "total_price": 317,
	//    "nm_id": 2389212,
	//    "brand": "Vivienne Sabo",
	//    "status": 202
	//  }
	//],
	//"locale": "en",
	//"internal_signature": "",
	//"customer_id": "test",
	//"delivery_service": "meest",
	//"shardkey": "9",
	//"sm_id": 99,
	//"date_created": "2021-11-26T06:22:19Z",
	//"oof_shard": "1"
	//}`))
	if err != nil {
		return err
	}
	fmt.Printf("publish: %+v\n", publish1)
	return nil
}
