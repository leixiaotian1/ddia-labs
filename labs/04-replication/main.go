package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== DDIA Labs: 数据复制 (Replication) ===")
	fmt.Println()

	demos := []struct {
		name        string
		description string
		path        string
	}{
		{
			name:        "主从复制 (Master-Slave)",
			description: "演示同步与异步复制及其一致性差异",
			path:        "master-slave",
		},
		{
			name:        "多主复制 (Multi-Master)",
			description: "演示并发写入冲突及 LWW 解决方法",
			path:        "multi-master",
		},
		{
			name:        "无主复制 (Leaderless)",
			description: "演示 Quorum (W+R>N) 读写机制",
			path:        "leaderless",
		},
	}

	fmt.Println("可用的复制模式演示:")
	fmt.Println()
	for i, demo := range demos {
		fmt.Printf("%d. %s\n", i+1, demo.name)
		fmt.Printf("   %s\n", demo.description)
		fmt.Printf("   路径: labs/04-replication/%s\n\n", demo.path)
	}

	fmt.Println("运行示例:")
	fmt.Println("  cd labs/04-replication/master-slave && go run main.go")
}

