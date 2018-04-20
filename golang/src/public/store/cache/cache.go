package cache;


import (
	"public/store/store"
)

type Cache = store.Store

func NewCache(config string) (Cache, error) {
	return store.NewStore(config)
}