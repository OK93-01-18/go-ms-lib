package controllers

import (
	"github.com/dgraph-io/ristretto"
	"time"
)

type Cacher interface {
	// Get returns the value (if any) and a boolean representing whether the
	// value was found or not. The value can be nil and the boolean can be true at
	// the same time.
	Get(interface{}) (interface{}, bool)

	// SetWithTTL works like Set but adds a key-value pair to the cache that will expire
	// after the specified TTL (time to live) has passed. A zero value means the value never
	// expires, which is identical to calling Set. A negative value is a no-op and the value
	// is discarded.
	SetWithTTL(interface{}, interface{}, int64, time.Duration) bool

	Wait()

	// Del deletes the key-value item from the cache if it exists.
	Del(interface{})

	// GetTTL returns the TTL for the specified key and a bool that is true if the
	// item was found and is not expired.
	GetTTL(interface{}) (time.Duration, bool)
}

type Config struct {
	// NumCounters determines the number of counters (keys) to keep that hold
	// access frequency information. It's generally a good idea to have more
	// counters than the max cache capacity, as this will improve eviction
	// accuracy and subsequent hit ratios.
	//
	// For example, if you expect your cache to hold 1,000,000 items when full,
	// NumCounters should be 10,000,000 (10x). Each counter takes up roughly
	// 3 bytes (4 bits for each counter * 4 copies plus about a byte per
	// counter for the bloom filter). Note that the number of counters is
	// internally rounded up to the nearest power of 2, so the space usage
	// may be a little larger than 3 bytes * NumCounters.
	NumCounters int64
	// MaxCost can be considered as the cache capacity, in whatever units you
	// choose to use.
	//
	// For example, if you want the cache to have a max capacity of 100MB, you
	// would set MaxCost to 100,000,000 and pass an item's number of bytes as
	// the `cost` parameter for calls to Set. If new items are accepted, the
	// eviction process will take care of making room for the new item and not
	// overflowing the MaxCost value.
	MaxCost int64
	// BufferItems determines the size of Get buffers.
	//
	// Unless you have a rare use case, using `64` as the BufferItems value
	// results in good performance.
	BufferItems int64
}

// Cache - cache for controller
type Cache struct {
	instance *ristretto.Cache
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	return c.instance.Get(key)
}

func (c *Cache) SetWithTTL(key, value interface{}, cost int64, ttl time.Duration) bool {
	return c.instance.SetWithTTL(key, value, cost, ttl)
}

func (c *Cache) Wait() {
	c.instance.Wait()
}

func (c *Cache) Del(key interface{}) {
	c.instance.Del(key)
}

func (c *Cache) GetTTL(key interface{}) (time.Duration, bool) {
	return c.instance.GetTTL(key)
}

func NewControllerCache(conf *Config) (Cacher, error) {
	rInstance, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: conf.NumCounters,
		MaxCost:     conf.MaxCost,
		BufferItems: conf.BufferItems,
	})
	if err != nil {
		return nil, err
	}

	return &Cache{
		instance: rInstance,
	}, nil
}
