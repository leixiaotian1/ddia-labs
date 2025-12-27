package main

import (
	"fmt"
	"time"
)

// 简化的性能对比测试
// 注意：这是概念演示，实际性能对比需要考虑更多因素

func benchmarkWrite(name string, writeFunc func(int, string), count int) time.Duration {
	start := time.Now()
	for i := 0; i < count; i++ {
		writeFunc(i, fmt.Sprintf("value%d", i))
	}
	return time.Since(start)
}

func benchmarkRead(name string, readFunc func(int) (string, bool), count int) time.Duration {
	start := time.Now()
	for i := 0; i < count; i++ {
		readFunc(i)
	}
	return time.Since(start)
}

func main() {
	fmt.Println("=== B-tree vs LSM-tree 性能对比 ===")
	fmt.Println()
	fmt.Println("注意：这是简化演示，实际性能受多种因素影响")
	fmt.Println()

	// 这里只是演示框架，实际需要导入btree和lsm的实现
	// 为了演示，我们使用模拟数据

	writeCount := 1000
	readCount := 1000

	fmt.Println("写入性能对比：")
	fmt.Printf("  写入 %d 条记录\n", writeCount)

	// 模拟B-tree写入（随机I/O，较慢）
	btreeWriteTime := time.Duration(writeCount*100) * time.Microsecond
	fmt.Printf("  B-tree:   %v (模拟：随机I/O，较慢)\n", btreeWriteTime)

	// 模拟LSM-tree写入（顺序I/O，较快）
	lsmWriteTime := time.Duration(writeCount*20) * time.Microsecond
	fmt.Printf("  LSM-tree: %v (模拟：顺序I/O，较快)\n", lsmWriteTime)
	fmt.Printf("  LSM-tree写入速度约为B-tree的 %.1f 倍\n", float64(btreeWriteTime)/float64(lsmWriteTime))

	fmt.Println("\n读取性能对比：")
	fmt.Printf("  读取 %d 条记录\n", readCount)

	// 模拟B-tree读取（直接定位，较快）
	btreeReadTime := time.Duration(readCount*10) * time.Microsecond
	fmt.Printf("  B-tree:   %v (模拟：直接定位，较快)\n", btreeReadTime)

	// 模拟LSM-tree读取（可能需要查多个SSTable，较慢）
	lsmReadTime := time.Duration(readCount*50) * time.Microsecond
	fmt.Printf("  LSM-tree: %v (模拟：可能查多个SSTable，较慢)\n", lsmReadTime)
	fmt.Printf("  B-tree读取速度约为LSM-tree的 %.1f 倍\n", float64(lsmReadTime)/float64(btreeReadTime))

	fmt.Println("\n=== 性能总结 ===")
	fmt.Println("B-tree:")
	fmt.Println("  ✓ 读取性能优秀（稳定，可预测）")
	fmt.Println("  ✗ 写入性能一般（随机I/O）")
	fmt.Println("  ✓ 适合OLTP系统，读多写少")
	fmt.Println("\nLSM-tree:")
	fmt.Println("  ✗ 读取性能一般（可能需要查多个层级）")
	fmt.Println("  ✓ 写入性能优秀（顺序I/O，高吞吐）")
	fmt.Println("  ✓ 适合写密集型应用，时序数据")
	fmt.Println("\n选择建议：")
	fmt.Println("- 读多写少 -> B-tree")
	fmt.Println("- 写多读少 -> LSM-tree")
	fmt.Println("- 需要范围查询 -> B-tree")
	fmt.Println("- 高写入吞吐量 -> LSM-tree")
}
