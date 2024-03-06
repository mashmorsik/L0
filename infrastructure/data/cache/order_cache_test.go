package cache

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/test/testdata"
	"testing"
	"time"
)

func TestOrderCache_Set_Get_success(t *testing.T) {
	log.BuildLogger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	orderCache := NewOrderCache(ctx, time.Hour)

	orderCache.Set(testdata.TestOrder1.OrderUid, *testdata.TestOrder1, time.Hour)

	cachedOrder, ok := orderCache.Get(testdata.TestOrder1.OrderUid)
	if !ok {
		t.Errorf("failed: order is not in cache")
	}
	assert.Equal(t, cachedOrder, testdata.TestOrder1)
}

func TestOrderCache_Set_Get_fail(t *testing.T) {
	log.BuildLogger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	orderCache := NewOrderCache(ctx, time.Hour)

	orderCache.Set(testdata.TestOrder1.OrderUid, *testdata.TestOrder1, time.Hour)

	cachedOrder, ok := orderCache.Get(testdata.TestOrder2.OrderUid)
	if ok {
		t.Errorf("failed: order is in cache")
	}
	assert.NotEqual(t, cachedOrder, testdata.TestOrder1)
}

func TestOrderCache_itemTypeCheck(t *testing.T) {
	log.BuildLogger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	orderCache := NewOrderCache(ctx, time.Hour)

	orderCache.Set(testdata.TestOrder1.OrderUid, *testdata.TestOrder1, time.Hour)

	item := &Item{
		Order:    *testdata.TestOrder1,
		Eviction: time.Now().Add(time.Second * 5),
	}

	_, ok := orderCache.isInvalidType(*item)
	if !ok {
		t.Errorf("failed: cached type error")
	}
}
