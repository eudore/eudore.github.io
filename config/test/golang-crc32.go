package main

import (
	"fmt"
	"strconv"
	"sync"
	"hash/crc32"
)

var keyBufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 256)
		return &b
	},
}

func crc32kv(key string) (uint32){
	bufp := keyBufPool.Get().(*[]byte)
	n := copy(*bufp, key)
	cs := crc32.ChecksumIEEE((*bufp)[:n])
	keyBufPool.Put(bufp)
	return cs
}

func main() {
	for i:=1;i<7;i++ {
		k := "key"+ strconv.Itoa(i)
		v := crc32kv(k)
		fmt.Printf("crc32 hash: k=%s v=%d\n",k,v)
    }
}

/*
crc32 hash: k=key1 v=744252496
crc32 hash: k=key2 v=3042260458
crc32 hash: k=key3 v=3260155260
crc32 hash: k=key4 v=1547079903
crc32 hash: k=key5 v=724672585
crc32 hash: k=key6 v=2990076403
*/