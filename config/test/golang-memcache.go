package main

import(
	"time"
	"fmt"
	"reflect"
	"github.com/pangudashu/memcache"
)

func main(){
	s1 := &memcache.Server{Address: "127.0.0.1:12001", Weight: 50}
	//s2 := &memcache.Server{Address: "127.0.0.1:12002", Weight: 50}
	mc, err := memcache.NewMemcache([]*memcache.Server{s1})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(reflect.TypeOf(mc))
	mc.SetRemoveBadServer(true)
	mc.SetTimeout(time.Second*2, time.Second, time.Second)

	//mc.Set("key1",false)
	//mc.Set("key2",1)
	//mc.Set("key3",222)
	//mc.Set("key4",4.63)
	mc.Set("123",15)
	//mc.Set("key6",222)
	fmt.Println(mc.Get("123"))
	fmt.Println(mc.Get("key1"))
	fmt.Println(mc.Get("key2"))
	fmt.Println(mc.Get("key3"))
	fmt.Println(mc.Get("key4"))
	fmt.Println(mc.Get("key5"))
	//mc.Delete("test_key")
	
	mc.Close()
}