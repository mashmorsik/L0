package models

type Order struct {
	OrderUid          string   `json:"order_uid" validate:"required"`
	TrackNumber       string   `json:"track_number" validate:"required"`
	Entry             string   `json:"entry" validate:"required"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale" validate:"required"`
	InternalSignature string   `json:"internal_signature" validate:"required"`
	CustomerId        string   `json:"customer_id" validate:"required"`
	DeliveryService   string   `json:"delivery_service" validate:"required"`
	Shardkey          string   `json:"shardkey" validate:"required"`
	SmId              int64    `json:"sm_id" validate:"required"`
	DateCreated       string   `json:"date_created" validate:"required"`
	OofShard          string   `json:"oof_shard" validate:"required"`
}
