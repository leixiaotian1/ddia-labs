package query

import (
	"fmt"
	"strings"
	"github.com/ddia-labs/labs/14-simple-db/storage"
	"github.com/ddia-labs/labs/14-simple-db/index"
	"github.com/ddia-labs/labs/14-simple-db/transaction"
)

// Engine 负责协调各个组件执行指令
type Engine struct {
	storage *storage.DiskStorage
	index   *index.Index
	lm      *transaction.LockManager
}

func NewEngine(s *storage.DiskStorage, i *index.Index, lm *transaction.LockManager) *Engine {
	return &Engine{
		storage: s,
		index:   i,
		lm:      lm,
	}
}

// Execute 模拟 SQL 解析和执行
// 支持指令: 
// - SET key value
// - GET key
func (e *Engine) Execute(command string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid command")
	}

	action := strings.ToUpper(parts[0])
	key := parts[1]

	switch action {
	case "SET":
		if len(parts) < 3 {
			return "", fmt.Errorf("SET requires a value")
		}
		value := parts[2]
		
		// 协调事务、存储和索引
		unlock := e.lm.LockKey(key)
		defer unlock()

		offset, err := e.storage.Write(key, value)
		if err != nil {
			return "", err
		}
		e.index.Put(key, offset)
		return "OK", nil

	case "GET":
		offset, ok := e.index.Get(key)
		if !ok {
			return "(nil)", nil
		}
		_, val, err := e.storage.ReadAt(offset)
		if err != nil {
			return "", err
		}
		return val, nil

	default:
		return "", fmt.Errorf("unknown command: %s", action)
	}
}

