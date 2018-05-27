package memcache

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"public/cache"
	"github.com/eudore/gomemcache/memcache"
)

// Store Memstore adapter.
type Store struct {
	conn     *memcache.Client
	conninfo []string
}

// NewMemStore create new memcache adapter.
func NewMemStore() cache.Cache {
	return &Store{}
}

// Get get value from memcache.
func (rc *Store) Get(key string) []byte {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil
			//return err
		}
	}
	if item, err := rc.conn.Get(key); err == nil {
		return item.Value
	}
	return nil
}

// GetMulti get value from memcache.
func (rc *Store) GetMulti(keys []string) [][]byte {
	var rv [][]byte
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil
		}
	}
	mv, err := rc.conn.GetMulti(keys)
	if err == nil {
		for _, v := range mv {
			rv = append(rv, v.Value)
		}
		return rv
	}
	return rv
}

// Put put value to memcache.
func (rc *Store) Put(key string, val []byte, timeout time.Duration) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	item := memcache.Item{Key: key,Value: val, Expiration: int32(timeout / time.Second)}
	return rc.conn.Set(&item)
}

// Delete delete value in memcache.
func (rc *Store) Delete(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.Delete(key)
}

// Incr increase counter.
func (rc *Store) Incr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Increment(key, 1)
	return err
}

// Decr decrease counter.
func (rc *Store) Decr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Decrement(key, 1)
	return err
}

// Get All keys
func (rc *Store) GetAllKeys() ([]string,error) {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil,err
		}
	}
	return rc.conn.GetAllKeys()

}

// Get keys size
func (rc *Store) Size() (int,error) {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return 0,err
		}
	}
	return rc.conn.Size()

}

// IsExist check value exists in memcache.
func (rc *Store) IsExist(key string) bool {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false
		}
	}
	_, err := rc.conn.Get(key)
	return !(err != nil)
}

// ClearAll clear all stored in memcache.
func (rc *Store) ClearAll() error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.FlushAll()
}

// StartAndGC start memcache adapter.
// config string is like {"conn":"connection info"}.
// if connecting error, return.
func (rc *Store) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	rc.conninfo = strings.Split(cf["conn"], ";")
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return nil
}

// connect to memcache and keep the connection.
func (rc *Store) connectInit() error {
	rc.conn = memcache.New(rc.conninfo...)
	return nil
}

func init() {
	cache.Register("memcache", NewMemStore)
}
