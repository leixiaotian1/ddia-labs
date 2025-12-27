package main

import (
	"fmt"
)

// GCounter is a Grow-only Counter (state-based CRDT)
type GCounter struct {
	ID     string
	counts map[string]int
}

func NewGCounter(id string) *GCounter {
	return &GCounter{
		ID:     id,
		counts: make(map[string]int),
	}
}

func (c *GCounter) Inc() {
	c.counts[c.ID]++
}

func (c *GCounter) Value() int {
	total := 0
	for _, v := range c.counts {
		total += v
	}
	return total
}

func (c *GCounter) Merge(other *GCounter) {
	for id, count := range other.counts {
		if count > c.counts[id] {
			c.counts[id] = count
		}
	}
}

func main() {
	fmt.Println("=== CRDT: G-Counter (增长计数器) 演示 ===")
	fmt.Println()
	fmt.Println("CRDT: 无冲突复制数据类型。允许在没有中央协调的情况下实现最终一致性。")
	fmt.Println()

	// 模拟两个节点 Node-1 和 Node-2
	node1 := NewGCounter("Node-1")
	node2 := NewGCounter("Node-2")

	fmt.Println("1. 两个节点各自独立执行加 1 操作")
	node1.Inc()
	node1.Inc() // Node-1 = 2
	node2.Inc() // Node-2 = 1

	fmt.Printf("   Node-1 当前感知值: %d\n", node1.Value())
	fmt.Printf("   Node-2 当前感知值: %d\n", node2.Value())

	fmt.Println("\n2. 执行双向 Merge (同步数据)")
	node1.Merge(node2)
	node2.Merge(node1)

	fmt.Printf("   Merge 后 Node-1 最终值: %d\n", node1.Value())
	fmt.Printf("   Merge 后 Node-2 最终值: %d\n", node2.Value())

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("1. G-Counter 只允许增加，每个节点维护自己的计数，Merge 时取各节点计数的最大值。")
	fmt.Println("2. CRDT 的核心是操作满足 结合律、交换律和幂等性 (ACI)。")
	fmt.Println("3. 这种机制使得分布式系统在网络分区情况下仍能保持高可用，并在网络恢复后自动达成一致。")
}

