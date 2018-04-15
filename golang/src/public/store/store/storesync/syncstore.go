package storesync;

import (
	"public/store/store"
)
type Sync struct {
	Master	*store.Store
	Slave 	[]store.Store
}

func (s *Sync) Get(key string) interface{} {
	return s.Master.Get(key)
}

func (s *Sync) GetMulti(keys []string) []interface{} {
	return s.Master.GetMulti(key)	
}

func (s *Sync) Put(key string, val interface{}, timeout time.Duration) error {
	if err := s.Master.Put(key,val,timeout);err != nil {
		for k := range Slave {
			k.Put(key,val,timeout)
		}
	}
	return err
}

// increase cached int value by key, as a counter.
func (s *Sync) Delete(key string) error {
	if err := s.Master.Delete(key);err != nil {
		for k := range Slave {
			k.Delete(key)
		}
	}	
	return err
}

func (s *Sync) Incr(key string) error {
	if err := s.Master.Incr(key);err != nil {
		for k := range Slave {
			k.Incr(key)
		}
	}	
	return err
}

// decrease cached int value by key, as a counter.
func (s *Sync) Decr(key string) error {
	if err := s.Master.StartAndGC(config);err != nil {
		for k := range Slave {
			k.StartAndGC(config)
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

// clear all cache.
func (s *Sync) ClearAll() error {
	if err := s.Master.ClearAll();err != nil {
		for k := range Slave {
			k.ClearAll()
		}
	}	
	return err
}

// start gc routine based on config string settings.
func (s *Sync) StartAndGC(config string) error {
	return nil
}


