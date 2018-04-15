package cache;

import (
	"testing"
	"time"
	"public/log"
	_ "public/store/store/memcache"
	"public/store/cache"
)

func TestCache(t *testing.T) {
	Token, err := cache.NewCache(`{"store": "memcache","name": "token","conn":"127.0.0.1:12003"}`)
	log.Info("global init token:",err)
	log.Info(Token.Put("key1", []byte("22"), 10 * time.Second))
	log.Info(Token.Get("key1"))
}