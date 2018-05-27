package cache

import (
	"sync"
	"time"
	"net/http"
	"encoding/json"

	kv "public/cache"
	"public/session"

	//"github.com/astaxie/beego/session"
	//"github.com/bradfitz/gomemcache/memcache"
)

var mempder = &MemProvider{}
var client kv.Cache

// SessionStore memcache session store
type SessionStore struct {
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime int64
}

// Set value in memcache session
func (rs *SessionStore) Set(key, value interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values[key] = value
	return nil
}

// Get value in memcache session
func (rs *SessionStore) Get(key interface{}) interface{} {
	rs.lock.RLock()
	defer rs.lock.RUnlock()
	if v, ok := rs.values[key]; ok {
		return v
	}
	return nil
}

// Delete value in memcache session
func (rs *SessionStore) Delete(key interface{}) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	delete(rs.values, key)
	return nil
}

// Flush clear all values in memcache session
func (rs *SessionStore) Flush() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	rs.values = make(map[interface{}]interface{})
	return nil
}

// SessionID get memcache session id
func (rs *SessionStore) SessionID() string {
	return rs.sid
}

// SessionRelease save session values to memcache
func (rs *SessionStore) SessionRelease(w http.ResponseWriter) {
	b, err := session.EncodeGob(rs.values)
	if err != nil {
		return
	}
	client.Put(rs.sid,b,time.Duration(rs.maxlifetime) * time.Second)
}

// MemProvider memcache session provider
type MemProvider struct {
	maxlifetime int64
	cachename 	string
	conninfo 	string
}

// SessionInit init memcache session
// savepath like
// e.g. 127.0.0.1:9090
func (rp *MemProvider) SessionInit(maxlifetime int64, config string) (err error) {
	var cf map[string]string
	err = json.Unmarshal([]byte(config), &cf)
	if err != nil {
		return
	}
	cacheName := cf["cache"]

	rp.cachename = cacheName
	rp.maxlifetime = maxlifetime
	rp.conninfo = config
	client,err = kv.NewCache(cacheName,config)
	return
}

// SessionRead read memcache session by sid
func (rp *MemProvider) SessionRead(sid string) (session.Store, error) {
	item := client.Get(sid)
	var kv map[interface{}]interface{}
	if len(item) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		var err error
		kv, err = session.DecodeGob(item)
		if err != nil {
			return nil, err
		}
	}
	rs := &SessionStore{sid: sid, values: kv, maxlifetime: rp.maxlifetime}
	return rs, nil
}

// SessionExist check memcache session exist by sid
func (rp *MemProvider) SessionExist(sid string) bool {
	if item := client.Get(sid); len(item) == 0 {
		return false
	}
	return true
}

// SessionRegenerate generate new sid for memcache session
func (rp *MemProvider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	var contain []byte
	if item := client.Get(sid); len(item) == 0 {
		client.Put(sid,[]byte(""),time.Duration(rp.maxlifetime) * time.Second)
	} else {
		client.Put(sid,item,time.Duration(rp.maxlifetime) * time.Second)
		contain = item
	}

	var kv map[interface{}]interface{}
	if len(contain) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		var err error
		kv, err = session.DecodeGob(contain)
		if err != nil {
			return nil, err
		}
	}

	rs := &SessionStore{sid: sid, values: kv, maxlifetime: rp.maxlifetime}
	return rs, nil
}

// SessionDestroy delete memcache session by id
func (rp *MemProvider) SessionDestroy(sid string) error {
	return client.Delete(sid)
}

// SessionGC Impelment method, no used.
func (rp *MemProvider) SessionGC() {
}

// SessionAll return all activeSession
func (rp *MemProvider) SessionAll() int {
	return 0
}

func init() {
	session.Register("cache", mempder)
}
