package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== DDIA Labs: 一致性模型 (Consistency Models) ===")
	fmt.Println()

	demos := []struct {
		name        string
		description string
		path        string
	}{
		{
			name:        "线性一致性 (Linearizability)",
			description: "演示最强一致性模型，原子生效的读写",
			path:        "linearizability",
		},
		{
			name:        "最终一致性 (Eventual Consistency)",
			description: "演示常见的弱一致性模型及其短暂不一致现象",
			path:        "eventual",
		},
		{
			name:        "因果一致性 (Causal Consistency)",
			description: "演示通过维护因果关系保证操作顺序",
			path:        "causal",
		},
	}

	fmt.Println("可用的一致性模型演示:")
	fmt.Println()
	for i, demo := range demos {
		fmt.Printf("%d. %s\n", i+1, demo.name)
		fmt.Printf("   %s\n", demo.description)
		fmt.Printf("   路径: labs/07-consistency/%s\n\n", demo.path)
	}

	fmt.Println("运行示例:")
	fmt.Println("  cd labs/07-consistency/linearizability && go run main.go")
}

