package cache;

import (
	"testing"
	"time"
	_ "public/cache/memcache"
	"public/cache"
)

func TestStore(t *testing.T) {
	bm, err := cache.NewCache("memcache",`{"cache": "memcache","conn":"127.0.0.1:12003"}`)	
	t.Log(err)
	t.Log(bm.Put("astax", []byte("22"), 10 * time.Second))
	t.Log(bm.Put("astaxie", []byte("sss"), 10 * time.Second))
	t.Log(bm.Get("astax"))
	t.Log(bm.GetMulti([]string{"astaxie","astax"}))
}
