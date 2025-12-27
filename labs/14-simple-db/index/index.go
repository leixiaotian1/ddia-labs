package index

// Index is a simple interface for indexing
type Index interface {
	Put(key string, pos int)
	Get(key string) (int, bool)
}

// MemoryIndex stores key and its position/reference
type MemoryIndex struct {
	data map[string]int
}

func NewMemoryIndex() *MemoryIndex {
	return &MemoryIndex{data: make(map[string]int)}
}

func (i *MemoryIndex) Put(key string, pos int) {
	i.data[key] = pos
}

func (i *MemoryIndex) Get(key string) (int, bool) {
	pos, ok := i.data[key]
	return pos, ok
}

