package cache;

import (
	"testing"
	"fmt"
	"time"
	_ "public/cache/memcache"
	"public/cache"
)

func TestStore(t *testing.T) {
	bm, err := cache.NewStore("memcache",`{"cache": "memcache","conn":"127.0.0.1:12003"}`)	
	fmt.Println(err)
	fmt.Println(bm.Put("astax", []byte("22"), 10 * time.Second))
	fmt.Println(bm.Put("astaxie", []byte("sss"), 10 * time.Second))
	fmt.Println(bm.Get("astax"))
	fmt.Println(bm.GetMulti([]string{"astaxie","astax"}))
}