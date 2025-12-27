package transaction

import (
	"sync"
)

// LockManager handles row-level concurrency control
type LockManager struct {
	locks sync.Map
}

func NewLockManager() *LockManager {
	return &LockManager{}
}

// LockKey ensures exclusive access to a specific key during write operations
func (lm *LockManager) LockKey(key string) func() {
	l, _ := lm.locks.LoadOrStore(key, &sync.Mutex{})
	mtx := l.(*sync.Mutex)
	mtx.Lock()
	// Return an unlock function for easy deferring
	return func() { mtx.Unlock() }
}
