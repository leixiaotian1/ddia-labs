package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== DDIA Labs: 索引结构 ===\n")

	demos := []struct {
		name        string
		description string
		path        string
	}{
		{
			name:        "哈希索引",
			description: "演示哈希索引的O(1)查找和等值查询",
			path:        "hash-index",
		},
		{
			name:        "B-tree索引",
			description: "演示B-tree索引的范围查询和有序遍历",
			path:        "btree-index",
		},
	}

	fmt.Println("可用的索引结构演示：\n")
	for i, demo := range demos {
		fmt.Printf("%d. %s\n", i+1, demo.name)
		fmt.Printf("   %s\n", demo.description)
		fmt.Printf("   路径: %s\n\n", demo.path)
	}

	fmt.Println("请进入对应的子目录运行具体的demo：")
	fmt.Println("  cd hash-index && go run main.go")
}

