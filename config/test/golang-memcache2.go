package main

import(
	"fmt"
    "time"
	"strconv"
	"reflect"
	"github.com/bradfitz/gomemcache/memcache"
)

func main(){
	//mc := memcache.New("127.0.0.1:12001")
	mc := memcache.New("127.0.0.1:12001", "127.0.0.1:12002", "127.0.0.1:12003")
    mc.Timeout = time.Second * 5
	fmt.Println(reflect.TypeOf(mc))
    mc.Set(&memcache.Item{Key: "key1", Value: []byte("String")})
    mc.Set(&memcache.Item{Key: "key2", Value: []byte("0"),Expiration: 3})
    mc.Set(&memcache.Item{Key: "key3", Value: []byte("15")})
    mc.Set(&memcache.Item{Key: "key4", Value: []byte("4.63")})
    mc.Set(&memcache.Item{Key: "key5", Value: []byte("this is string")})
    mc.Set(&memcache.Item{Key: "key6", Value: []byte("一个字符串")})

    for i:=1;i<7;i++ {
        key := "key"+ strconv.Itoa(i)
    	it, err := mc.Get(key)
    	if err == nil{
			fmt.Printf("Fetch key:%s data:%s %d\n",key,string(it.Value),it.Expiration )
    	}

    }
}