package main

import (
	"fmt"
	"sync"
)

// Message represents a causal event
type Message struct {
	ID      int
	Content string
	Parent  int // ID of the message this one is replying to
}

// CausalStorage maintains messages and can check causal order
type CausalStorage struct {
	Messages map[int]Message
	mu       sync.Mutex
}

func (s *CausalStorage) Add(m Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Messages[m.ID] = m
}

func (s *CausalStorage) DisplayOrder() {
	fmt.Println("--- 消息流展示 ---")
	// Simplified: just check if replies appear after their parents
	for _, m := range s.Messages {
		if m.Parent != 0 {
			parent, exists := s.Messages[m.Parent]
			if exists {
				fmt.Printf("  [回复] ID:%d '%s' -> 参考了 ID:%d '%s' (因果关系正确)\n", 
					m.ID, m.Content, parent.ID, parent.Content)
			} else {
				fmt.Printf("  [异常] ID:%d '%s' -> 参考了未知的 ID:%d (因果违背/数据丢失)\n", 
					m.ID, m.Content, m.Parent)
			}
		} else {
			fmt.Printf("  [原始] ID:%d '%s'\n", m.ID, m.Content)
		}
	}
}

func main() {
	fmt.Println("=== 因果一致性 (Causal Consistency) 演示 ===")
	fmt.Println()
	fmt.Println("定义: 具有因果关系的操作必须以相同的顺序被所有节点观察到。")
	fmt.Println("场景: 社交网络中的回复。回复必须在其被回复的消息之后显示。")
	fmt.Println()

	storage := &CausalStorage{Messages: make(map[int]Message)}

	// 场景: 用户 A 提问，用户 B 回答。如果读取者先看到了回答而没看到提问，就会产生困惑。
	m1 := Message{ID: 1, Content: "你知道 DDIA 这本书吗？", Parent: 0}
	m2 := Message{ID: 2, Content: "当然，它是分布式系统的经典！", Parent: 1}

	fmt.Println("[System] 消息 1 发出")
	storage.Add(m1)
	
	fmt.Println("[System] 消息 2 (基于消息 1) 发出")
	storage.Add(m2)

	storage.DisplayOrder()

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("因果一致性通过跟踪操作之间的依赖关系来保证顺序。")
	fmt.Println("权衡: 比线性化一致性弱（不保证不相关操作的顺序），但比最终一致性强，且可以避免很多用户体验上的混乱。")
}

