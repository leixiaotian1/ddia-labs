package main

import (
	"fmt"
	"sync"
	"time"
)

// Data represents a simple key-value pair
type Data struct {
	Value   string
	Version int
}

// Node represents a database node (Master or Slave)
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

func (n *Node) Write(key string, value string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	version := 0
	if existing, ok := n.Storage[key]; ok {
		version = existing.Version + 1
	}
	n.Storage[key] = Data{Value: value, Version: version}
}

func (n *Node) Read(key string) (Data, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	data, ok := n.Storage[key]
	return data, ok
}

// MasterSlaveSystem simulates a master-slave replication system
type MasterSlaveSystem struct {
	Master   *Node
	Slaves   []*Node
	syncMode bool // true: synchronous, false: asynchronous
}

func NewMasterSlaveSystem(slaveCount int, syncMode bool) *MasterSlaveSystem {
	slaves := make([]*Node, slaveCount)
	for i := 0; i < slaveCount; i++ {
		slaves[i] = NewNode(fmt.Sprintf("Slave-%d", i))
	}
	return &MasterSlaveSystem{
		Master:   NewNode("Master"),
		Slaves:   slaves,
		syncMode: syncMode,
	}
}

func (s *MasterSlaveSystem) Replicate(key string, value string) {
	// 1. Write to Master
	s.Master.Write(key, value)
	fmt.Printf("[Master] 写入成功: %s = %s\n", key, value)

	if s.syncMode {
		// Synchronous replication: wait for all slaves to acknowledge
		var wg sync.WaitGroup
		for _, slave := range s.Slaves {
			wg.Add(1)
			go func(sl *Node) {
				defer wg.Done()
				// Simulate network latency
				time.Sleep(100 * time.Millisecond)
				sl.Write(key, value)
				fmt.Printf("[%s] 同步复制成功\n", sl.ID)
			}(slave)
		}
		wg.Wait()
		fmt.Println("[System] 所有从节点已确认，写入完成")
	} else {
		// Asynchronous replication: return immediately after master write
		for _, slave := range s.Slaves {
			go func(sl *Node) {
				// Simulate network latency
				time.Sleep(200 * time.Millisecond)
				sl.Write(key, value)
				fmt.Printf("[%s] 异步复制完成 (后台任务)\n", sl.ID)
			}(slave)
		}
		fmt.Println("[System] 异步复制已触发，主节点已返回")
	}
}

func main() {
	fmt.Println("=== 主从复制 (Master-Slave Replication) 演示 ===")
	fmt.Println()

	// 1. 演示同步复制
	fmt.Println("--- 场景 1: 同步复制 (Synchronous) ---")
	fmt.Println("特点: 必须等待从节点确认写入，一致性强，但延迟高。")
	syncSystem := NewMasterSlaveSystem(2, true)
	syncSystem.Replicate("user1", "Alice")

	val, _ := syncSystem.Slaves[0].Read("user1")
	fmt.Printf("[Slave-0] 立即读取 user1: %s\n", val.Value)
	fmt.Println()

	// 2. 演示异步复制
	fmt.Println("--- 场景 2: 异步复制 (Asynchronous) ---")
	fmt.Println("特点: 不等待从节点确认，延迟低，但可能存在读取不一致（复制延迟）。")
	asyncSystem := NewMasterSlaveSystem(2, false)
	asyncSystem.Replicate("user2", "Bob")

	// 立即读取，可能读不到
	val, ok := asyncSystem.Slaves[0].Read("user2")
	if !ok {
		fmt.Println("[Slave-0] 立即读取 user2: 未找到 (存在复制延迟！)")
	} else {
		fmt.Printf("[Slave-0] 立即读取 user2: %s\n", val.Value)
	}

	// 等待一段时间再读
	fmt.Println("[System] 等待 300ms...")
	time.Sleep(300 * time.Millisecond)
	val, _ = asyncSystem.Slaves[0].Read("user2")
	fmt.Printf("[Slave-0] 延迟后读取 user2: %s\n", val.Value)

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("主从复制适用于读多写少的场景。")
	fmt.Println("同步复制 vs 异步复制的权衡:")
	fmt.Println("- 同步: 强一致性，低可用性 (一个从节点挂了写操作就阻塞)")
	fmt.Println("- 异步: 高可用性，弱一致性 (可能存在读取不一致)")
}
