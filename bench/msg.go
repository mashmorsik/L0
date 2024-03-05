package bench

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/mashmorsik/L0/pkg/models"
	"strconv"
	"strings"
	"time"
)

func GenerateMsg() *models.Order {
	f := gofakeit.New(0)

	orderID := genID(f)
	trackNumber := strings.ReplaceAll(strings.ToUpper(f.HipsterWord()), " ", "") +
		strings.ReplaceAll(strings.ToUpper(f.HipsterWord()), " ", "")
	items := generateItems(f, trackNumber)
	goodsTotal := func() int {
		var t int
		for _, item := range items {
			t += item.TotalPrice
		}
		return t
	}()
	deliveryCost := f.RandomInt([]int{100, 200, 300, 400})
	customFee := f.IntRange(0, 200)

	fakeOrd := models.Order{
		OrderUid:    orderID,
		TrackNumber: trackNumber,
		Entry:       trackNumber[:3],
		Delivery: models.Delivery{
			Name:    f.Name(),
			Phone:   f.Phone(),
			Zip:     f.Zip(),
			City:    f.City(),
			Address: f.Address().Address,
			Region:  f.RandomString([]string{"Moscow Region", "Spb", "Ekat"}),
			Email:   f.Email(),
		},
		Payment: models.Payment{
			Transaction:  orderID,
			RequestId:    genID(f),
			Currency:     f.CurrencyShort(),
			Provider:     f.RandomString([]string{"applepay", "paypal", "onlinekassa"}),
			Amount:       deliveryCost + goodsTotal + customFee,
			PaymentDt:    f.PastDate().Unix(),
			Bank:         f.RandomString([]string{"sber", "tinkoff", "alpha", "raif"}),
			DeliveryCost: deliveryCost,
			GoodsTotal:   goodsTotal,
			CustomFee:    customFee,
		},
		Items:             items,
		Locale:            f.RandomString([]string{"en", "ru"}),
		InternalSignature: f.Word(),
		CustomerId:        genID(f),
		DeliveryService:   f.Company(),
		Shardkey:          strconv.Itoa(f.IntRange(1, 100)),
		SmId:              int64(f.IntRange(1, 100)),
		DateCreated:       time.Now().Format(time.RFC3339),
		OofShard:          strconv.Itoa(f.IntRange(1, 100)),
	}

	return &fakeOrd
}

func generateItems(f *gofakeit.Faker, trackNum string) []models.Item {
	itemsNum := f.IntRange(1, 3)
	items := make([]models.Item, itemsNum)

	for i := 0; i < itemsNum; i++ {
		price := int(f.Product().Price)
		sale := f.IntRange(5, 99)
		count := f.IntRange(1, 5)

		items[i] = models.Item{
			ChrtId:      f.Int(),
			TrackNumber: trackNum,
			Price:       price,
			Rid:         genID(f),
			Name:        f.ProductName(),
			Sale:        sale,
			Size:        strconv.Itoa(f.IntRange(1, 1000)),
			Count:       count,
			TotalPrice:  (price - (price / 100 * sale)) * count,
			NmId:        f.IntRange(1, 10000000),
			Brand:       f.Company(),
			Status:      f.HTTPStatusCode(),
		}
	}

	return items
}

func genID(f *gofakeit.Faker) string {
	hash := md5.New()
	hash.Write([]byte(f.Adverb()))
	hashedBytes := hash.Sum(nil)
	ID := hex.EncodeToString(hashedBytes)

	return ID
}
