package main

import (
	"fmt"
	"time"
)

// LWWRegister is a Last-Write-Wins Register
type LWWRegister struct {
	Value     string
	Timestamp int64
}

func (r *LWWRegister) Merge(other LWWRegister) {
	if other.Timestamp > r.Timestamp {
		r.Value = other.Value
		r.Timestamp = other.Timestamp
	}
}

func main() {
	fmt.Println("=== CRDT: LWW-Register (最后写入获胜寄存器) 演示 ===")
	fmt.Println()
	fmt.Println("特点: 冲突解决策略为保留时间戳最新的数据。")
	fmt.Println()

	reg1 := LWWRegister{Value: "Version-1", Timestamp: time.Now().UnixNano()}
	
	// 模拟稍后的写入
	time.Sleep(10 * time.Millisecond)
	reg2 := LWWRegister{Value: "Version-2", Timestamp: time.Now().UnixNano()}

	fmt.Printf("Register 1: %s (ts: %d)\n", reg1.Value, reg1.Timestamp)
	fmt.Printf("Register 2: %s (ts: %d)\n", reg2.Value, reg2.Timestamp)

	fmt.Println("\n执行 Merge (1 <- 2)...")
	reg1.Merge(reg2)

	fmt.Printf("最终结果: %s\n", reg1.Value)

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("1. LWW 是一种简单但有效的冲突处理方法。")
	fmt.Println("2. 严重依赖时钟同步。如果节点时钟偏差较大，可能会丢失较晚但时戳较小的数据。")
}

