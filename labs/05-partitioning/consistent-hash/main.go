package main

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strconv"
)

// Node represents a physical server
type Node struct {
	ID   string
	Data map[string]string
}

// ConsistentHash manages nodes and mapping
type ConsistentHash struct {
	Nodes    map[string]*Node
	Ring     []uint32
	NodeMap  map[uint32]string // Hash to Node ID
	Replicas int              // Virtual nodes for better distribution
}

func NewConsistentHash(replicas int) *ConsistentHash {
	return &ConsistentHash{
		Nodes:    make(map[string]*Node),
		NodeMap:  make(map[uint32]string),
		Replicas: replicas,
	}
}

func hash(key string) uint32 {
	h := sha1.New()
	h.Write([]byte(key))
	b := h.Sum(nil)
	// Return first 4 bytes as uint32
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

func (ch *ConsistentHash) AddNode(nodeID string) {
	ch.Nodes[nodeID] = &Node{ID: nodeID, Data: make(map[string]string)}
	for i := 0; i < ch.Replicas; i++ {
		// Virtual node hash
		h := hash(nodeID + "#" + strconv.Itoa(i))
		ch.Ring = append(ch.Ring, h)
		ch.NodeMap[h] = nodeID
	}
	sort.Slice(ch.Ring, func(i, j int) bool { return ch.Ring[i] < ch.Ring[j] })
	fmt.Printf("[System] 添加节点 %s (%d 个虚拟节点)\n", nodeID, ch.Replicas)
}

func (ch *ConsistentHash) RemoveNode(nodeID string) {
	delete(ch.Nodes, nodeID)
	newRing := []uint32{}
	for _, h := range ch.Ring {
		if ch.NodeMap[h] != nodeID {
			newRing = append(newRing, h)
		} else {
			delete(ch.NodeMap, h)
		}
	}
	ch.Ring = newRing
	fmt.Printf("[System] 移除节点 %s\n", nodeID)
}

func (ch *ConsistentHash) GetNode(key string) string {
	if len(ch.Ring) == 0 {
		return ""
	}
	h := hash(key)
	// Find the first virtual node hash >= h
	idx := sort.Search(len(ch.Ring), func(i int) bool { return ch.Ring[i] >= h })
	if idx == len(ch.Ring) {
		idx = 0 // Wrap around to the start of the ring
	}
	return ch.NodeMap[ch.Ring[idx]]
}

func (ch *ConsistentHash) Write(key, value string) {
	nodeID := ch.GetNode(key)
	ch.Nodes[nodeID].Data[key] = value
	fmt.Printf("[Write] Key='%s' (hash=%d) -> Node %s\n", key, hash(key), nodeID)
}

func main() {
	fmt.Println("=== 一致性哈希 (Consistent Hashing) 演示 ===")
	fmt.Println()
	fmt.Println("特点: 节点增减时，只有一小部分数据需要迁移。常用于分布式缓存和动态分区。")
	fmt.Println()

	ch := NewConsistentHash(3) // 3个虚拟节点

	ch.AddNode("Server-A")
	ch.AddNode("Server-B")
	ch.AddNode("Server-C")

	fmt.Println("\n--- 初始写入 ---")
	keys := []string{"user-1", "user-2", "user-3", "user-4", "user-5"}
	for _, k := range keys {
		ch.Write(k, "data_"+k)
	}

	fmt.Println("\n--- 动态扩缩容演示 ---")
	// 记录 user-3 原本在哪个节点
	oldNode := ch.GetNode("user-3")
	fmt.Printf("user-3 原本在: %s\n", oldNode)

	ch.RemoveNode("Server-B")
	newNode := ch.GetNode("user-3")
	fmt.Printf("移除 Server-B 后，user-3 路由到: %s\n", newNode)

	fmt.Println("\n=== 总结 ===")
	fmt.Println("一致性哈希通过虚拟节点解决了数据分布不均的问题，并极大地降低了节点变更时的重新映射开销。")
	fmt.Println("它是分布式系统实现弹性伸缩的基础。")
}

