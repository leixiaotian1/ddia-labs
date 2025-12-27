package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== DDIA Labs: 共识算法 (Consensus) ===")
	fmt.Println()

	demos := []struct {
		name        string
		description string
		path        string
	}{
		{
			name:        "Raft 共识算法",
			description: "演示 Leader 选举和高可用故障恢复",
			path:        "raft",
		},
		{
			name:        "Basic Paxos 简化演示",
			description: "演示 Prepare 和 Accept 的投票逻辑",
			path:        "paxos-simple",
		},
	}

	fmt.Println("可用的共识算法演示:")
	fmt.Println()
	for i, demo := range demos {
		fmt.Printf("%d. %s\n", i+1, demo.name)
		fmt.Printf("   %s\n", demo.description)
		fmt.Printf("   路径: labs/08-consensus/%s\n\n", demo.path)
	}

	fmt.Println("运行示例:")
	fmt.Println("  cd labs/08-consensus/raft && go run main.go")
}

