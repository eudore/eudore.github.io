package main

import(
	"fmt"
	"time"
	"strconv"
	"reflect"
	"gomemcache/memcache"
)

func main(){
	//mc := memcache.New("127.0.0.1:12001")
	mc := memcache.New("127.0.0.1:12001","127.0.0.1:12002","127.0.0.1:12003")
	mc.Timeout = time.Second * 2
	mc.SetHash(jenkins)
	fmt.Println(reflect.TypeOf(mc))
	mc.Set(&memcache.Item{Key: "key1", Value: []byte("String"),Expiration: 300})
	mc.Set(&memcache.Item{Key: "key2", Value: []byte("0"),Expiration: 3})
	mc.Set(&memcache.Item{Key: "key3", Value: []byte("15"),Expiration: 300})
	mc.Set(&memcache.Item{Key: "key4", Value: []byte("4.63"),Expiration: 300})
	mc.Set(&memcache.Item{Key: "key5", Value: []byte("this is string"),Expiration: 300})
	mc.Set(&memcache.Item{Key: "key6", Value: []byte("一个字符串"),Expiration: 300})
	mc.Set(&memcache.Item{Key: "ksss", Value: []byte("一个字符串"),Expiration: 300})
	
	fmt.Println(mc.Size())
	fmt.Println(mc.GetAllKeys())
	
	for i:=1;i<8;i++ {
		key := "key"+ strconv.Itoa(i)
		it, err := mc.Get(key)
		if err == nil{
			fmt.Printf("Fetch key:%s data:%s %d\n",key,string(it.Value),it.Expiration )
		}
	}
}


func jenkins(key []byte) (uint32) {
	var hash uint32 = 0
	for _, b := range key {
		hash += uint32(b)
		hash += (hash << 10)
		hash ^= (hash >> 6)
	}
	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)
	return hash
}
