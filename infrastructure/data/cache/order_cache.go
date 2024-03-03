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
	Ctx              context.Context
	EvictionDuration time.Duration
}

type Item struct {
	Order    models.Order
	Eviction time.Time
}

func NewCache() *OrderCache {
	return &OrderCache{}
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

	item, ok := c.itemTypeCheck(foundItem)
	if !ok {
		return nil, false
	}

	return &item.Order, true
}

func (c *OrderCache) evictionWorker() {
	ticker := time.NewTicker(c.EvictionDuration)
	for {
		select {
		case <-c.Ctx.Done():
			return
		case <-ticker.C:
			cache.Range(func(key any, value any) bool {
				item, ok := c.itemTypeCheck(value)
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

func (c *OrderCache) itemTypeCheck(foundItem any) (item *Item, invalid bool) {
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
