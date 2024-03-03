package models

import (
	"math/big"
	"time"
)

type Payment struct {
	Transaction  string    `json:"transaction"`
	RequestId    string    `json:"request_id"`
	Currency     string    `json:"currency"`
	Provider     string    `json:"provider"`
	Amount       big.Rat   `json:"amount"`
	PaymentDt    time.Time `json:"payment_dt"`
	Bank         string    `json:"bank"`
	DeliveryCost big.Rat   `json:"delivery_cost"`
	GoodsTotal   big.Rat   `json:"goods_total"`
	CustomFee    big.Rat   `json:"custom_fee"`
}
