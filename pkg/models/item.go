package models

import "github.com/shopspring/decimal"

type Item struct {
	ChrtId      int             `json:"chrt_id"`
	TrackNumber string          `json:"track_number"`
	Price       decimal.Decimal `json:"price"`
	Rid         string          `json:"rid"`
	Name        string          `json:"name"`
	Sale        int             `json:"sale"`
	Size        string          `json:"size"`
	Count       int             `json:"count"`
	TotalPrice  decimal.Decimal `json:"total_price"`
	NmId        int             `json:"nm_id"`
	Brand       string          `json:"brand"`
	Status      int             `json:"status"`
}
