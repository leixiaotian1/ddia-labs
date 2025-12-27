package main

import (
	"fmt"
	"sync"
	"time"
)

// Node represents a distributed node with its own view of data
type Node struct {
	ID      string
	Storage map[string]string
	mu      sync.RWMutex
}

func NewNode(id string) *Node {
	return &Node{ID: id, Storage: make(map[string]string)}
}

func (n *Node) Write(key, val string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Storage[key] = val
}

func (n *Node) Read(key string) string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Storage[key]
}

func main() {
	fmt.Println("=== 最终一致性 (Eventual Consistency) 演示 ===")
	fmt.Println()
	fmt.Println("特点: 允许短时间内的不一致，只要没有新的更新，最终所有副本都会达到一致。")
	fmt.Println()

	nodeA := NewNode("Node-A")
	nodeB := NewNode("Node-B")

	fmt.Println("[User] 在 Node-A 写入 key='status', value='Active'")
	nodeA.Write("status", "Active")

	// 模拟复制延迟
	fmt.Println("[System] 模拟复制延迟 (异步同步中...)")
	
	fmt.Printf("[Reader] 在 Node-A 读取: %s\n", nodeA.Read("status"))
	fmt.Printf("[Reader] 在 Node-B 读取: %s (不一致！还没有同步过来)\n", nodeB.Read("status"))

	fmt.Println("[System] 等待 200ms 同步完成...")
	time.Sleep(200 * time.Millisecond)
	
	// 后台同步逻辑
	nodeB.Write("status", nodeA.Read("status"))

	fmt.Printf("[Reader] 在 Node-A 读取: %s\n", nodeA.Read("status"))
	fmt.Printf("[Reader] 在 Node-B 读取: %s (最终一致了)\n", nodeB.Read("status"))

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("最终一致性是目前大多数分布式数据库（如 DNS, DynamoDB）采用的模型，因为它对延迟和可用性最友好。")
	fmt.Println("挑战: 开发者必须处理这种短暂的不一致（如“读你所写”一致性的缺失）。")
}

