package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== DDIA Labs: 存储引擎对比 ===\n")

	demos := []struct {
		name        string
		description string
		path        string
	}{
		{
			name:        "B-tree存储引擎",
			description: "演示B-tree的存储和检索机制",
			path:        "btree",
		},
		{
			name:        "LSM-tree存储引擎",
			description: "演示LSM-tree的追加写入和合并机制",
			path:        "lsm",
		},
		{
			name:        "性能对比",
			description: "对比B-tree和LSM-tree的性能特征",
			path:        "benchmark",
		},
	}

	fmt.Println("可用的存储引擎演示：\n")
	for i, demo := range demos {
		fmt.Printf("%d. %s\n", i+1, demo.name)
		fmt.Printf("   %s\n", demo.description)
		fmt.Printf("   路径: %s\n\n", demo.path)
	}

	fmt.Println("请进入对应的子目录运行具体的demo：")
	fmt.Println("  cd btree && go run main.go")
}

