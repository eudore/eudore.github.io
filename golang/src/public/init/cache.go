package init;

import (
	"public/cache"
)

func init() {
	go initcache()
}

func initcache() {
	bm,err := cache.NewCache("memcache",`{"conn":"127.0.0.1:12001"}`)
	if(err==nil){
		bm.Put("weer/public","file",8640000)
	}
}