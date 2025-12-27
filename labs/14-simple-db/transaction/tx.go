package transaction

import (
	"sync"
)

// TransactionManager handles simple transactions
type TransactionManager struct {
	locks sync.Map // Simple row-level locking simulation
}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{}
}

func (tm *TransactionManager) Lock(key string) {
	lock, _ := tm.locks.LoadOrStore(key, &sync.Mutex{})
	lock.(*sync.Mutex).Lock()
}

func (tm *TransactionManager) Unlock(key string) {
	lock, ok := tm.locks.Load(key)
	if ok {
		lock.(*sync.Mutex).Unlock()
	}
}

