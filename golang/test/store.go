package main;

import (
	"fmt"
	"time"
	_ "public/store/store/memcache"
	"public/store/store"
)

func main() {
	bm, err := store.NewStore(`{"store": "memcache","conn":"127.0.0.1:12003"}`)	
	fmt.Println(err)
	fmt.Println(bm.Put("astax", []byte("22"), 10 * time.Second))
	fmt.Println(bm.Put("astaxie", []byte("sss"), 10 * time.Second))
	fmt.Println(bm.Get("astax"))
	fmt.Println(bm.GetMulti([]string{"astaxie","astax"}))
}