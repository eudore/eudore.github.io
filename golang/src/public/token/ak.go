package token

import (
	"sync"
)

type Ak interface {
	AkGet(string) []byte
	AkSet(string, []byte)
	AkDel(string)
}

type defaultAk struct {
	kv		map[string][]byte
	lock	sync.RWMutex
}

func (ak *defaultAk) AkGet(key string) []byte {
	ak.lock.Lock()
	defer ak.lock.Unlock()
	if v, ok := ak.kv[key]; ok {
		return v
	}
	return nil
}

func (ak *defaultAk) AkSet(key string, value []byte) {
	ak.lock.Lock()
	defer ak.lock.Unlock()
	ak.kv[key] = value
}

func (ak *defaultAk) AkDel(key string) {
	ak.lock.Lock()
	defer ak.lock.Unlock()
	delete(ak.kv, key)
}

var dfak Ak = &defaultAk{	kv:	make(map[string][]byte),}

func AkGet(key string) {
	dfak.AkGet(key)
}

func AkSet(key string, value []byte) {
	dfak.AkSet(key, value)
}

func AkDel(key string) {
	dfak.AkDel(key)
}

func SetAk(ak Ak) {
	dfak = ak
}