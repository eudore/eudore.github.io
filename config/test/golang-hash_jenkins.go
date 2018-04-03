package main


import (
	"fmt"
	"strconv"
)

func jenkins_one_at_a_time_hash(key []byte) (uint32) {
	var hash uint32
	hash =0

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


func main() {
	for i:=1;i<7;i++ {
		k := "key"+ strconv.Itoa(i)
		v := jenkins_one_at_a_time_hash([]byte(k))
		fmt.Printf("jenkins hash: k=%s v=%d\n",k,v)
    }
}