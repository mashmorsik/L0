package order

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/mashmorsik/L0/infrastructure/data/cache"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/pkg/models"
	"github.com/mashmorsik/L0/test/testdata"
	"github.com/nats-io/nats.go"
	"reflect"
	"testing"
	"time"
)

func TestCreateOrder_CreateNewOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := testdata.NewMockRepository(ctrl)
	mockRepo.EXPECT().CreateOrder(*testdata.TestOrder1).Return(nil)

	ctx := context.Background()
	orderCache := cache.NewOrderCache(ctx, time.Hour)

	createOrder := NewCreateOrder(mockRepo, orderCache)

	err := createOrder.CreateNewOrder(*testdata.TestOrder1)
	if err != nil {
		t.Errorf("can't create new order, err: %s", err)
	}

	cachedOrder, ok := orderCache.Get(testdata.TestOrder1.OrderUid)
	if !ok {
		t.Errorf("failed: can't get created order from cache, orderID: %s\n", testdata.TestOrder1.OrderUid)
	}

	assert.Equal(t, testdata.TestOrder1, *cachedOrder)
}

func TestCreateOrder_GetOrderFromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := testdata.NewMockRepository(ctrl)

	ctx := context.Background()
	orderCache := cache.NewOrderCache(ctx, time.Hour)

	orderCache.Set(testdata.TestOrder1.OrderUid, *testdata.TestOrder1, time.Hour)
	orderCache.Set(testdata.TestOrder2.OrderUid, *testdata.TestOrder2, time.Hour)

	createOrder := NewCreateOrder(mockRepo, orderCache)

	cachedOrder, err := createOrder.GetOrderFromCache(testdata.TestOrder2.OrderUid)
	if err != nil {
		t.Errorf("failed: can't get created order from cache, orderID: %s\n", testdata.TestOrder2.OrderUid)
	}

	assert.Equal(t, testdata.TestOrder2, *cachedOrder)
}

func TestCreateOrder_LoadCache(t *testing.T) {
	log.BuildLogger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx := &sql.Tx{}

	mockRepo := testdata.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetOrders().Return([]*models.Order{testdata.TestOrder1, testdata.TestOrder2}, nil)
	mockRepo.EXPECT().GetOrdersIDTx(gomock.Any()).Return([]string{"acceae55b8caf53492ba722354d92af6", "ff563feb7b2b84b6test"}, nil)

	ctx := context.Background()
	orderCache := cache.NewOrderCache(ctx, time.Hour)

	createOrder := NewCreateOrder(mockRepo, orderCache)
	err := createOrder.LoadCache()
	if err != nil {
		t.Errorf("failed: can't load cache, err: %s\n", err)
	}

	ordersID, err := mockRepo.GetOrdersIDTx(tx)
	if err != nil {
		t.Errorf("failed: can't get ordersID, err: %s\n", err)
	}

	for _, id := range ordersID {
		_, err = createOrder.GetOrderFromCache(id)
		if err != nil {
			t.Errorf("failed: order is not in cache, err: %s\n", err)
		}
	}
}

func TestUnmarshalOrder(t *testing.T) {
	type args struct {
		msg nats.Msg
	}

	jsonData := testdata.JSONData

	tests := []struct {
		name    string
		args    args
		want    *models.Order
		wantErr error
	}{
		{name: "unmarshal_filled_order_successfully",
			args: args{msg: nats.Msg{
				Data: []byte(jsonData),
			}}, want: testdata.TestOrder1, wantErr: nil},
		{name: "unmarshal_empty_order_successfully",
			args: args{msg: nats.Msg{
				Data: []byte(`{"jsonData":"empty"}`),
			}}, want: &models.Order{}, wantErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalOrder(tt.args.msg)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UnmarshalOrder() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalOrder() got = %v, want %v\n", got, tt.want)
			}
		})
	}
}
