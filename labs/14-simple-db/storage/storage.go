package storage

import (
	"fmt"
	"sort"
)

type Entry struct {
	Key   string
	Value string
}

// LSMStorage is a simplified LSM-tree storage engine
type LSMStorage struct {
	memTable   map[string]string
	ssTables   [][]Entry
	threshold  int
}

func NewLSMStorage(threshold int) *LSMStorage {
	return &LSMStorage{
		memTable:  make(map[string]string),
		ssTables:  make([][]Entry, 0),
		threshold: threshold,
	}
}

func (s *LSMStorage) Put(key, value string) {
	s.memTable[key] = value
	if len(s.memTable) >= s.threshold {
		s.flush()
	}
}

func (s *LSMStorage) Get(key string) (string, bool) {
	if val, ok := s.memTable[key]; ok {
		return val, true
	}
	// Check SSTables from newest to oldest
	for i := len(s.ssTables) - 1; i >= 0; i-- {
		idx := sort.Search(len(s.ssTables[i]), func(j int) bool {
			return s.ssTables[i][j].Key >= key
		})
		if idx < len(s.ssTables[i]) && s.ssTables[i][idx].Key == key {
			return s.ssTables[i][idx].Value, true
		}
	}
	return "", false
}

func (s *LSMStorage) flush() {
	fmt.Println("[Storage] MemTable 达到阈值，刷新到 SSTable...")
	entries := make([]Entry, 0, len(s.memTable))
	for k, v := range s.memTable {
		entries = append(entries, Entry{Key: k, Value: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	s.ssTables = append(s.ssTables, entries)
	s.memTable = make(map[string]string)
}

