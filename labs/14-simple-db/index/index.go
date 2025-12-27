package index

import "sync"

// Index maps key to disk offset
type Index struct {
	mu    sync.RWMutex
	table map[string]int64
}

func NewIndex() *Index {
	return &Index{table: make(map[string]int64)}
}

func (i *Index) Put(key string, offset int64) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.table[key] = offset
}

func (i *Index) Get(key string) (int64, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	offset, ok := i.table[key]
	return offset, ok
}
