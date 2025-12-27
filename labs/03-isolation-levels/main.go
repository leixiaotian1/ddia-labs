package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== DDIA Labs: 事务隔离级别测试 ===")
	fmt.Println()

	demos := []struct {
		name        string
		description string
		path        string
	}{
		{
			name:        "脏读 (Dirty Read)",
			description: "演示Read Uncommitted级别下的脏读问题",
			path:        "dirty-read",
		},
		{
			name:        "不可重复读 (Non-repeatable Read)",
			description: "演示Read Committed级别下的不可重复读问题",
			path:        "non-repeatable",
		},
		{
			name:        "幻读 (Phantom Read)",
			description: "演示Repeatable Read级别下的幻读问题",
			path:        "phantom-read",
		},
		{
			name:        "可串行化 (Serializable)",
			description: "演示Serializable级别如何避免所有并发问题",
			path:        "serializable",
		},
	}

	fmt.Println("可用的隔离级别演示：")
	fmt.Println()
	for i, demo := range demos {
		fmt.Printf("%d. %s\n", i+1, demo.name)
		fmt.Printf("   %s\n", demo.description)
		fmt.Printf("   路径: %s\n\n", demo.path)
	}

	fmt.Println("请进入对应的子目录运行具体的demo：")
	fmt.Println("  cd dirty-read && go run main.go")
}
