package models

import (
	"github.com/shopspring/decimal"
)

type Payment struct {
	Transaction  string          `json:"transaction"`
	RequestId    string          `json:"request_id"`
	Currency     string          `json:"currency"`
	Provider     string          `json:"provider"`
	Amount       decimal.Decimal `json:"amount"`
	PaymentDt    int64           `json:"payment_dt"`
	Bank         string          `json:"bank"`
	DeliveryCost decimal.Decimal `json:"delivery_cost"`
	GoodsTotal   decimal.Decimal `json:"goods_total"`
	CustomFee    decimal.Decimal `json:"custom_fee"`
}
