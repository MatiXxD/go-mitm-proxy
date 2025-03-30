package proxy

import (
	"sync"
)

type MemProxyRepository struct {
	cache map[string]map[string][]byte
	mu    sync.RWMutex
}

func NewMemProxyRepository() *MemProxyRepository {
	return &MemProxyRepository{
		cache: make(map[string]map[string][]byte),
	}
}

func (pr *MemProxyRepository) Store(k string, v map[string][]byte) {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	pr.cache[k] = v
}

func (pr *MemProxyRepository) Get(k string) (map[string][]byte, bool) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	val, ok := pr.cache[k]
	if !ok {
		return nil, false
	}

	return val, true
}
