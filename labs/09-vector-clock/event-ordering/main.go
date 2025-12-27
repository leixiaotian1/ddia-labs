package main

import (
	"fmt"
)

// VectorClock represents the clock state of a node
type VectorClock map[string]int

func NewVectorClock() VectorClock {
	return make(VectorClock)
}

// Tick increments the counter for the local node
func (vc VectorClock) Tick(nodeID string) {
	vc[nodeID]++
}

// Merge combines another vector clock into this one (takes the max of each component)
func (vc VectorClock) Merge(other VectorClock) {
	for nodeID, count := range other {
		if count > vc[nodeID] {
			vc[nodeID] = count
		}
	}
}

// Compare checks the causal relationship between two clocks
func (vc VectorClock) Compare(other VectorClock) string {
	less := false
	greater := false

	// Union of all node IDs in both clocks
	allNodes := make(map[string]bool)
	for n := range vc {
		allNodes[n] = true
	}
	for n := range other {
		allNodes[n] = true
	}

	for n := range allNodes {
		v1 := vc[n]
		v2 := other[n]
		if v1 < v2 {
			less = true
		} else if v1 > v2 {
			greater = true
		}
	}

	if less && !greater {
		return "Before"
	} else if !less && greater {
		return "After"
	} else if !less && !greater {
		return "Identical"
	}
	return "Concurrent"
}

func main() {
	fmt.Println("=== 向量时钟 (Vector Clock) 演示 ===")
	fmt.Println()
	fmt.Println("特点: 用于捕捉分布式系统中的因果关系 (Causal Relationship)。")
	fmt.Println()

	// 模拟三个节点 A, B, C
	vcA := NewVectorClock()
	vcB := NewVectorClock()
	vcC := NewVectorClock()

	fmt.Println("1. 节点 A 发生一个本地事件")
	vcA.Tick("A")
	fmt.Printf("   Clock A: %v\n", vcA)

	fmt.Println("\n2. 节点 A 发送消息给 B，B 接收并处理")
	// A 发送时的状态
	msgClock := make(VectorClock)
	for k, v := range vcA { msgClock[k] = v }
	
	vcB.Tick("B")
	vcB.Merge(msgClock)
	fmt.Printf("   Clock B: %v (继承了 A 的状态)\n", vcB)

	fmt.Println("\n3. 节点 C 发生一个独立事件")
	vcC.Tick("C")
	fmt.Printf("   Clock C: %v\n", vcC)

	fmt.Println("\n--- 因果关系分析 ---")
	fmt.Printf("A vs B: %s (A 发生在 B 之前)\n", vcA.Compare(vcB))
	fmt.Printf("B vs A: %s (B 发生在 A 之后)\n", vcB.Compare(vcA))
	fmt.Printf("B vs C: %s (B 和 C 是并发发生的，没有因果关系)\n", vcB.Compare(vcC))

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("1. 向量时钟通过维护一个包含所有节点计数器的向量来跟踪因果关系。")
	fmt.Println("2. 如果 v1[i] <= v2[i] 对所有 i 都成立，且至少有一个 j 满足 v1[j] < v2[j]，则 v1 发生在 v2 之前。")
	fmt.Println("3. 如果无法比较（互有大小），则这两个事件是并发的 (Concurrent)。")
}

