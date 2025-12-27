package main

import (
	"fmt"
	"strings"
)

// MapReduce 简化模拟
func main() {
	fmt.Println("=== MapReduce 简化版演示 ===")
	fmt.Println()
	fmt.Println("场景: 统计单词出现的频率 (Word Count)。")
	fmt.Println()

	input := []string{
		"hello world",
		"hello ddia",
		"world of distributed systems",
	}

	fmt.Println("--- 1. Map 阶段 ---")
	intermediate := make(map[string][]int)
	for _, doc := range input {
		words := strings.Fields(doc)
		for _, w := range words {
			intermediate[w] = append(intermediate[w], 1)
			fmt.Printf("   Map: %s -> 1\n", w)
		}
	}

	fmt.Println("\n--- 2. Reduce 阶段 ---")
	finalCount := make(map[string]int)
	for word, counts := range intermediate {
		sum := 0
		for _, c := range counts {
			sum += c
		}
		finalCount[word] = sum
		fmt.Printf("   Reduce: %s -> %d\n", word, sum)
	}

	fmt.Println()
	fmt.Println("--- 最终结果 ---")
	for k, v := range finalCount {
		fmt.Printf("%s: %d\n", k, v)
	}

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("1. Map: 并行处理数据，生成中间键值对。")
	fmt.Println("2. Reduce: 聚合中间结果。")
	fmt.Println("3. MapReduce 的核心价值在于将海量数据的计算分摊到廉价的计算集群上，并具备自动容错能力。")
}

