package main

import (
	"fmt"
	"github.com/ddia-labs/labs/14-simple-db/storage"
	"github.com/ddia-labs/labs/14-simple-db/transaction"
	"time"
)

// SimpleDB combines all components
type SimpleDB struct {
	storage *storage.LSMStorage
	tm      *transaction.TransactionManager
}

func NewSimpleDB() *SimpleDB {
	return &SimpleDB{
		storage: storage.NewLSMStorage(3),
		tm:      transaction.NewTransactionManager(),
	}
}

func (db *SimpleDB) Set(key, value string) {
	db.tm.Lock(key)
	defer db.tm.Unlock(key)
	db.storage.Put(key, value)
}

func (db *SimpleDB) Get(key string) (string, bool) {
	return db.storage.Get(key)
}

func main() {
	fmt.Println("=== SimpleDB: 综合数据库演示 ===")
	fmt.Println()
	fmt.Println("该示例整合了:")
	fmt.Println("1. 存储层: LSM-tree 结构")
	fmt.Println("2. 事务层: 行级锁模拟")
	fmt.Println("3. 接口层: 简单的 Get/Set 接口")
	fmt.Println()

	db := NewSimpleDB()

	// 演示写入
	fmt.Println("--- 写入测试 ---")
	db.Set("user:1", "Alice")
	db.Set("user:2", "Bob")
	db.Set("user:3", "Charlie")
	db.Set("user:4", "David") // 触发 LSM Flush

	// 演示并发写入和锁
	fmt.Println("\n--- 并发事务测试 ---")
	go func() {
		db.tm.Lock("user:1")
		fmt.Println("[Tx 1] 获取 user:1 锁，开始耗时操作...")
		time.Sleep(200 * time.Millisecond)
		db.storage.Put("user:1", "Alice-Revised")
		db.tm.Unlock("user:1")
		fmt.Println("[Tx 1] 提交并释放锁")
	}()

	time.Sleep(50 * time.Millisecond)
	fmt.Println("[Main] 尝试更新 user:1 (将被 Tx 1 阻塞)")
	db.Set("user:1", "Alice-Main")
	fmt.Println("[Main] 更新 user:1 成功")

	// 演示读取
	fmt.Println("\n--- 读取测试 ---")
	val, _ := db.Get("user:1")
	fmt.Printf("user:1 = %s\n", val)
	val2, _ := db.Get("user:4")
	fmt.Printf("user:4 = %s (从 SSTable 读取)\n", val2)

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("1. SimpleDB 展示了如何将存储引擎、事务管理整合在一起。")
	fmt.Println("2. LSM-tree 保证了高效写入，通过多层级查询保证读取。")
	fmt.Println("3. 事务管理器保证了并发访问的正确性。")
}

