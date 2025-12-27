package main

import (
	"fmt"
	"os"
	"github.com/ddia-labs/labs/14-simple-db/storage"
	"github.com/ddia-labs/labs/14-simple-db/index"
	"github.com/ddia-labs/labs/14-simple-db/transaction"
	"github.com/ddia-labs/labs/14-simple-db/query"
)

func main() {
	dbFile := "simple.db"
	defer os.Remove(dbFile)

	// 1. 初始化底层组件
	s, _ := storage.NewDiskStorage(dbFile)
	idx := index.NewIndex()
	lm := transaction.NewLockManager()
	defer s.Close()

	// 2. 初始化查询引擎 (Query Layer)
	engine := query.NewEngine(s, idx, lm)

	fmt.Println("=== SimpleDB: 完整架构演示 (含查询层) ===")
	fmt.Println("架构层级: Query(解析) -> Transaction(并发控制) -> Storage(磁盘) & Index(内存)")
	fmt.Println()

	// 演示 3: 通过查询层执行指令
	fmt.Println("--- 场景: 通过查询层执行指令 ---")
	
	commands := []string{
		"SET user:1 Alice",
		"SET user:2 Bob",
		"GET user:1",
		"SET user:1 Alice_Updated",
		"GET user:1",
		"GET user:999", // 不存在的 key
	}

	for _, cmd := range commands {
		result, err := engine.Execute(cmd)
		if err != nil {
			fmt.Printf("执行失败 [%s]: %v\n", cmd, err)
		} else {
			fmt.Printf("执行指令 [%s] -> 结果: %s\n", cmd, result)
		}
	}

	fmt.Println()
	fmt.Println("=== 为什么需要 Query 层？ ===")
	fmt.Println("1. 抽象细节：用户不需要知道磁盘 Offset 或如何加锁，只需发送字符串指令。")
	fmt.Println("2. 统一入口：所有并发冲突和组件协调都在 Query 层完成，保证了系统的正确性。")
	fmt.Println("3. 可扩展性：如果以后要支持 SQL 或 JSON 查询，只需在 Query 层增加解析逻辑。")
}
