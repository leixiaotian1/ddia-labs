package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== DDIA Labs: 数据分区 (Partitioning) ===")
	fmt.Println()

	demos := []struct {
		name        string
		description string
		path        string
	}{
		{
			name:        "范围分区 (Range Partitioning)",
			description: "根据 Key 的字典序范围划分，支持范围查询",
			path:        "range-partition",
		},
		{
			name:        "哈希分区 (Hash Partitioning)",
			description: "使用哈希散列 Key，保证负载均衡",
			path:        "hash-partition",
		},
		{
			name:        "一致性哈希 (Consistent Hashing)",
			description: "动态分区策略，支持节点平滑增减",
			path:        "consistent-hash",
		},
	}

	fmt.Println("可用的分区模式演示:")
	fmt.Println()
	for i, demo := range demos {
		fmt.Printf("%d. %s\n", i+1, demo.name)
		fmt.Printf("   %s\n", demo.description)
		fmt.Printf("   路径: labs/05-partitioning/%s\n\n", demo.path)
	}

	fmt.Println("运行示例:")
	fmt.Println("  cd labs/05-partitioning/range-partition && go run main.go")
}

