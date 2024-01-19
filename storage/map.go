package storage

import "sync"

type StoreMap struct {
	lock sync.RWMutex
	m    map[int64]Store
}

func NewStoreMap() *StoreMap {
	return &StoreMap{
		m: make(map[int64]Store),
	}
}

// get store or return nil
func (s *StoreMap) GetStore(chainID int64) Store {
	s.lock.RLock()
	defer s.lock.RUnlock()

	v := s.m[chainID]
	return v
}

func (s *StoreMap) SetStore(chainID int64, store Store) bool {
	if _, ok := store.(*PostgresStore); !ok {
		return false
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.m[chainID] = store

	return true
}

func (s *StoreMap) GetAll() []Store {
	s.lock.RLock()
	defer s.lock.RUnlock()

	stores := make([]Store, 0, len(s.m))
	for _, v := range s.m {
		stores = append(stores, v)
	}

	return stores
}
