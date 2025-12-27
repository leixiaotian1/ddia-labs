package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== DDIA Labs: 分布式事务 (Distributed Transactions) ===")
	fmt.Println()

	demos := []struct {
		name        string
		description string
		path        string
	}{
		{
			name:        "两阶段提交 (2PC)",
			description: "演示强原子性的阻塞式事务协议",
			path:        "two-phase-commit",
		},
		{
			name:        "Saga 补偿模式",
			description: "演示基于补偿操作的最终一致性长事务",
			path:        "saga",
		},
	}

	fmt.Println("可用的分布式事务演示:")
	fmt.Println()
	for i, demo := range demos {
		fmt.Printf("%d. %s\n", i+1, demo.name)
		fmt.Printf("   %s\n", demo.description)
		fmt.Printf("   路径: labs/06-distributed-tx/%s\n\n", demo.path)
	}

	fmt.Println("运行示例:")
	fmt.Println("  cd labs/06-distributed-tx/saga && go run main.go")
}

