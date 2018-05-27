package sync;

import (
	"time"
	"public/cache"
)
type Sync struct {
	Master	cache.Cache
	Slave 	[]cache.Cache
}

func (s *Sync) Get(key string) []byte {
	return s.Master.Get(key)
}

func (s *Sync) GetMulti(keys []string) [][]byte {
	return s.Master.GetMulti(keys)	
}

func (s *Sync) Put(key string, val []byte, timeout time.Duration) error {
	err := s.Master.Put(key,val,timeout)
	if err == nil {
		for _,k := range s.Slave {
			k.Put(key,val,timeout)
		}
	}
	return err
}

// increase cached int value by key, as a counter.
func (s *Sync) Delete(key string) error {
	err := s.Master.Delete(key)
	if err == nil {
		for _,k := range s.Slave {
			k.Delete(key)
		}
	}	
	return err
}

func (s *Sync) Incr(key string) error {
	err := s.Master.Incr(key)
	if err == nil {
		for _,k := range s.Slave {
			k.Incr(key)
		}
	}	
	return err
}

// decrease cached int value by key, as a counter.
func (s *Sync) Decr(key string) error {
	err := s.Master.Decr(key)
	if err == nil {
		for _,k := range s.Slave {
			k.Decr(key)
		}
	}	
	return err
}

// check if cached value exists or not.
func (s *Sync) IsExist(key string) bool {
	return s.Master.IsExist(key)
}

// get all keys
func (s *Sync) GetAllKeys() ([]string,error){
	return s.Master.GetAllKeys()	
}

// get all keys
func (s *Sync) Size() (int,error){
	return s.Master.Size()	
}

// clear all cache.
func (s *Sync) ClearAll() error {
	err := s.Master.ClearAll()
	if err == nil {
		for _, k := range s.Slave {
			k.ClearAll()
		}
	}	
	return err
}

// start gc routine based on config string settings.
func (s *Sync) StartAndGC(config string) error {
	return nil
}

func NewMemStore() cache.Cache {
	return &Sync{}
}

func init() {
	cache.Register("sync", NewMemStore)
}