package cache

import (
	"context"
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/pkg/models"
	"sync"
	"time"
)

var cache sync.Map

type OrderCache struct {
	Ctx                    context.Context
	evictionWorkerDuration time.Duration
}

func NewOrderCache(ctx context.Context, evictionWorkerDuration time.Duration) OrderCache {
	return OrderCache{Ctx: ctx, evictionWorkerDuration: evictionWorkerDuration}
}

type Item struct {
	Order    models.Order
	Eviction time.Time
}

func (c *OrderCache) Set(key string, order models.Order, expiration time.Duration) {
	cache.Store(key, Item{
		Order:    order,
		Eviction: time.Now().Add(expiration),
	})
}

func (c *OrderCache) Get(key string) (*models.Order, bool) {
	foundItem, ok := cache.Load(key)
	if !ok {
		return nil, false
	}

	item, ok := c.isInvalidType(foundItem)
	if !ok {
		return nil, false
	}

	return &item.Order, true
}

// evictionWorker iterates over the cached items every c.evictionWorkerDuration, checks if they are of the type *Item
// and deletes the orders whose storage time has expired.
func (c *OrderCache) evictionWorker() {
	ticker := time.NewTicker(c.evictionWorkerDuration)
	for {
		select {
		case <-c.Ctx.Done():
			return
		case <-ticker.C:
			cache.Range(func(key any, value any) bool {
				item, ok := c.isInvalidType(value)
				if !ok {
					return true
				}

				if item.Eviction.Before(time.Now()) {
					cache.Delete(key)
				}

				return true
			})
		}
	}
}

// isInvalidType returns false if the item is of the type *Item.
func (c *OrderCache) isInvalidType(foundItem any) (item *Item, invalid bool) {
	var cacheItem Item
	switch foundItem.(type) {
	case Item:
		cacheItem = foundItem.(Item)
	default:
		log.Errf("invalid cache item %#+v", foundItem)
		return nil, false
	}

	return &cacheItem, true
}
