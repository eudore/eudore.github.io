package store

import (
	"fmt"
	"time"
	"encoding/json"
	//_ "public/store/store/memcache"
)

// Cache interface contains all behaviors for cache adapter.
// usage:
//	cache.Register("file",cache.NewFileCache) // this operation is run in init method of file.go.
//	c,err := cache.NewCache("file","{....}")
//	c.Put("key",value, 3600 * time.Second)
//	v := c.Get("key")
//
//	c.Incr("counter")  // now is 1
//	c.Incr("counter")  // now is 2
//	count := c.Get("counter").(int)
type Store interface {
	// get cached value by key.
	Get(key string) []byte
	// GetMulti is a batch version of Get.
	GetMulti(keys []string) [][]byte
	// set cached value with key and expire time.
	Put(key string, val []byte, timeout time.Duration) error
	// delete cached value by key.
	Delete(key string) error
	// increase cached int value by key, as a counter.
	Incr(key string) error
	// decrease cached int value by key, as a counter.
	Decr(key string) error
	// check if cached value exists or not.
	IsExist(key string) bool
	// get all keys
	GetAllKeys() ([]string,error)
	// clear all cache.
	ClearAll() error
	// start gc routine based on config string settings.
	StartAndGC(config string) error
}

// Instance is a function create a new Cache Instance
type Instance func() Store

var adapters = make(map[string]Instance)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("store: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("store: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewCache Create a new cache driver by adapter name and config string.
// config need to be correct JSON as string: {"interval":360}.
// it will start gc automatically.
func NewStore(config string) (adapter Store, err error) {
	var cf map[string]string
	err = json.Unmarshal([]byte(config), &cf)
	if err != nil {
		adapter = nil
		return
	}
	adapterName := cf["store"]
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("store: unknown store adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
