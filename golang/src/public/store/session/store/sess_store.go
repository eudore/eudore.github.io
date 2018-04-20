package store

import (
	"net/http"
	"sync"
	"time"

	"public/store/store"
	"public/store/session"

	//"github.com/astaxie/beego/session"
	//"github.com/bradfitz/gomemcache/memcache"
)

var mempder = &MemProvider{}
var client store.Store

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
	//item := memcache.Item{Key: rs.sid, Value: b, Expiration: int32(rs.maxlifetime)}
	//client.Set(&item)
}

// MemProvider memcache session provider
type MemProvider struct {
	maxlifetime int64
	//conninfo    []string
	conninfo 	string
	poolsize    int
	password    string
}

// SessionInit init memcache session
// savepath like
// e.g. 127.0.0.1:9090
func (rp *MemProvider) SessionInit(maxlifetime int64, config string) (err error) {
	rp.maxlifetime = maxlifetime
	//rp.conninfo = strings.Split(savePath, ";")
	rp.conninfo=config
	client,err = store.NewStore(config)
	//client = memcache.New(rp.conninfo...)
	return
}

// SessionRead read memcache session by sid
func (rp *MemProvider) SessionRead(sid string) (session.Store, error) {
	if client == nil {
		if err := rp.connectInit(); err != nil {
			return nil, err
		}
	}
/*	item, err := client.Get(sid)
	if err != nil && err == memcache.ErrCacheMiss {
		rs := &SessionStore{sid: sid, values: make(map[interface{}]interface{}), maxlifetime: rp.maxlifetime}
		return rs, nil
	}*/
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
	if client == nil {
		if err := rp.connectInit(); err != nil {
			return false
		}
	}
	if item := client.Get(sid); len(item) == 0 {
		return false
	}
	return true
}

// SessionRegenerate generate new sid for memcache session
func (rp *MemProvider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	if client == nil {
		if err := rp.connectInit(); err != nil {
			return nil, err
		}
	}
	var contain []byte
	if item := client.Get(sid); len(item) == 0 {
		// oldsid doesn't exists, set the new sid directly
		// ignore error here, since if it return error
		// the existed value will be 0
/*		item.Key = sid
		item.Value = []byte("")
		item.Expiration = int32(rp.maxlifetime)*/
		client.Put(sid,[]byte(""),time.Duration(rp.maxlifetime) * time.Second)
	} else {
/*		client.Delete(oldsid)
		item.Key = sid
		item.Expiration = int32(rp.maxlifetime)
		client.Set(item)
		contain = item.Value*/
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
	if client == nil {
		if err := rp.connectInit(); err != nil {
			return err
		}
	}

	return client.Delete(sid)
}

func (rp *MemProvider) connectInit() (err error) {
	//client = memcache.New(rp.conninfo...)
	client,err = store.NewStore(rp.conninfo)
	return 
}

// SessionGC Impelment method, no used.
func (rp *MemProvider) SessionGC() {
}

// SessionAll return all activeSession
func (rp *MemProvider) SessionAll() int {
	return 0
}

func init() {
	session.Register("memcache", mempder)
}
