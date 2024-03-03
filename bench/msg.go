package bench

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mashmorsik/L0/pkg/models"
	"math/big"
)

func generateMsg() *models.Order {
	f := gofakeit.New(0)

	order := models.Order{
		OrderUid:    "",
		TrackNumber: "",
		Entry:       "",
		Delivery: models.Delivery{
			Name:    f.Name(),
			Phone:   f.Phone(),
			Zip:     f.Zip(),
			City:    f.City(),
			Address: f.Address().Address,
			Region:  "",
			Email:   f.Email(),
		},
		Payment: models.Payment{
			Transaction:  "",
			RequestId:    "",
			Currency:     f.CurrencyShort(),
			Provider:     "",
			Amount:       big.Rat{},
			PaymentDt:    f.PastDate(),
			Bank:         "",
			DeliveryCost: big.Rat{},
			GoodsTotal:   big.Rat{},
			CustomFee:    big.Rat{},
		},
		Items:             nil,
		Locale:            "",
		InternalSignature: "",
		CustomerId:        "",
		DeliveryService:   "",
		Shardkey:          "",
		SmId:              0,
		DateCreated:       f.PastDate(),
		OofShard:          "",
	}

	return &order
}
