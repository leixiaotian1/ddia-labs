package main

import (
	"fmt"
	"sync"
	"time"
)

// Data represents values with version/timestamp for conflict resolution
type Data struct {
	Value     string
	Timestamp int64
}

// Node represents a Master node in multi-master setup
type Node struct {
	ID      string
	Storage map[string]Data
	mu      sync.RWMutex
}

func NewNode(id string) *Node {
	return &Node{
		ID:      id,
		Storage: make(map[string]Data),
	}
}

// Write local write
func (n *Node) Write(key string, value string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Storage[key] = Data{
		Value:     value,
		Timestamp: time.Now().UnixNano(),
	}
}

// ReplicateFrom simulates receiving replication data from another master
func (n *Node) ReplicateFrom(key string, incoming Data) {
	n.mu.Lock()
	defer n.mu.Unlock()

	existing, exists := n.Storage[key]
	if !exists {
		n.Storage[key] = incoming
		fmt.Printf("[%s] 接受来自外部的写入: %s = %s\n", n.ID, key, incoming.Value)
		return
	}

	// Conflict resolution: Last Write Wins (LWW)
	if incoming.Timestamp > existing.Timestamp {
		n.Storage[key] = incoming
		fmt.Printf("[%s] 冲突解决 (LWW): 接受新值 %s (新时戳 %d > 旧时戳 %d)\n", 
			n.ID, incoming.Value, incoming.Timestamp, existing.Timestamp)
	} else {
		fmt.Printf("[%s] 冲突解决 (LWW): 保留旧值 %s (旧时戳 %d > 新时戳 %d)\n", 
			n.ID, existing.Value, existing.Timestamp, incoming.Timestamp)
	}
}

func main() {
	fmt.Println("=== 多主复制 (Multi-Master Replication) 演示 ===")
	fmt.Println()
	fmt.Println("特点: 多个节点都可以接受写入。适用于跨数据中心或离线操作场景。")
	fmt.Println("核心问题: 并发写入冲突。本示例演示 LWW (Last Write Wins) 冲突解决方法。")
	fmt.Println()

	masterA := NewNode("DC-Shanghai")
	masterB := NewNode("DC-London")

	// 1. 模拟在不同数据中心几乎同时写入
	fmt.Println("--- 场景: 并发冲突 ---")
	
	fmt.Println("[DC-Shanghai] 用户更新名字为 'Alice-New'")
	masterA.Write("username", "Alice-New")
	
	// 稍微延迟一点点，模拟 London 稍微晚一点的写入
	time.Sleep(10 * time.Millisecond)
	fmt.Println("[DC-London] 用户更新名字为 'Alice-Final'")
	masterB.Write("username", "Alice-Final")

	fmt.Println()
	fmt.Println("--- 触发异步复制同步数据 ---")
	
	// 获取 A 的数据同步给 B
	dataA := masterA.Storage["username"]
	masterB.ReplicateFrom("username", dataA)

	// 获取 B 的数据同步给 A
	dataB := masterB.Storage["username"]
	masterA.ReplicateFrom("username", dataB)

	fmt.Println()
	fmt.Println("--- 最终状态检查 ---")
	finalA, _ := masterA.Storage["username"]
	finalB, _ := masterB.Storage["username"]
	fmt.Printf("[DC-Shanghai] 最终值: %s\n", finalA.Value)
	fmt.Printf("[DC-London]   最终值: %s\n", finalB.Value)

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("多主复制提高了可用性和写入吞吐量。")
	fmt.Println("权衡:")
	fmt.Println("- 优势: 容忍单个数据中心故障，网络分区下仍可工作。")
	fmt.Println("- 挑战: 冲突解决 (Conflict Resolution) 非常复杂。")
	fmt.Println("- 常用方法: LWW, CRDT (无冲突复制数据类型), 提示修正等。")
}

