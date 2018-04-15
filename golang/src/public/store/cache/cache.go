package cache;


import (
	"fmt"
	"time"
	"encoding/json"
	"public/store/store"
)



type Cache interface {
	//NewCache(config string) (Cache,error)
	Get(key string) ([]byte,error)
	Put(key string,val []byte, timeout time.Duration) error
	Delete(key string) error
	IsExist(key string) bool
	Destory() error
}


// Instance is a function create a new Cache Instance
type Instance func(config string) (Cache,error)

var instances = make(map[string]Cache)
var adapters = make(map[string]Instance)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewCache Create a new cache driver by adapter name and config string.
// config need to be correct JSON as string: {"interval":360}.
// it will start gc automatically.
func NewCache(config string) (adapter Cache, err error) {
	var cf map[string]string
	err = json.Unmarshal([]byte(config), &cf)
	if err != nil {
		adapter = nil
		return
	}
	adapterName := cf["type"]
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("queue: unknown queue type name %q (forgot to import?)", adapterName)
		instanceFunc = NewDefaultCache
		//return
	}
	adapter,err = instanceFunc(config)
	if err != nil {
		adapter = nil
	}else{
		instances[cf["name"]] = adapter	
	}
	return
}



type defaultCache struct {
	store 	store.Store
	config 	string
}

func NewDefaultCache(config string) (c Cache,err error) {
	s,err := store.NewStore(config)
	if err != nil {
		return nil,err
	}
	c = &defaultCache{
		store: s,
		config: config,
	}
	return
}


func (c *defaultCache) Get(key string) ([]byte,error){
	return c.store.Get(key),nil
}

func (c *defaultCache)Put(key string,val []byte, timeout time.Duration) error {
	return c.store.Put(key,val,timeout)
}

func (c *defaultCache) Delete(key string) error {
	return c.store.Delete(key)
}

func (c *defaultCache) IsExist(key string) bool {
	return c.store.IsExist(key)
}

func (c *defaultCache) Destory() error {
	return c.store.ClearAll()
}
