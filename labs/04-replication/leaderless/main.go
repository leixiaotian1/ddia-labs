package main

import (
	"fmt"
	"sync"
	"time"
)

// Data represents values in leaderless system
type Data struct {
	Value     string
	Timestamp int64
}

// Node represents a database node in leaderless system
type Node struct {
	ID      string
	Storage map[string]Data
	Online  bool
	mu      sync.RWMutex
}

func NewNode(id string) *Node {
	return &Node{
		ID:      id,
		Storage: make(map[string]Data),
		Online:  true,
	}
}

// Cluster represents a leaderless cluster
type Cluster struct {
	Nodes []*Node
}

func NewCluster(count int) *Cluster {
	nodes := make([]*Node, count)
	for i := 0; i < count; i++ {
		nodes[i] = NewNode(fmt.Sprintf("Node-%d", i))
	}
	return &Cluster{Nodes: nodes}
}

// Quorum Write: Write to multiple nodes, succeed if W nodes confirm
func (c *Cluster) Write(key string, value string, W int) bool {
	successCount := 0
	data := Data{
		Value:     value,
		Timestamp: time.Now().UnixNano(),
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	fmt.Printf("[Client] 尝试并发写入 %d 个节点 (W=%d): %s = %s\n", len(c.Nodes), W, key, value)

	for _, node := range c.Nodes {
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			n.mu.Lock()
			defer n.mu.Unlock()
			
			if n.Online {
				n.Storage[key] = data
				mu.Lock()
				successCount++
				mu.Unlock()
				fmt.Printf("  [%s] 写入成功\n", n.ID)
			} else {
				fmt.Printf("  [%s] 离线，写入失败\n", n.ID)
			}
		}(node)
	}
	wg.Wait()

	if successCount >= W {
		fmt.Printf("[Client] 写入成功 (确认数 %d >= W=%d)\n", successCount, W)
		return true
	}
	fmt.Printf("[Client] 写入失败 (确认数 %d < W=%d)\n", successCount, W)
	return false
}

// Quorum Read: Read from multiple nodes, succeed if R nodes respond
func (c *Cluster) Read(key string, R int) (Data, bool) {
	var responses []Data
	var mu sync.Mutex
	var wg sync.WaitGroup

	fmt.Printf("[Client] 尝试并发读取 %d 个节点 (R=%d)\n", len(c.Nodes), R)

	for _, node := range c.Nodes {
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			n.mu.RLock()
			defer n.mu.RUnlock()

			if n.Online {
				if val, ok := n.Storage[key]; ok {
					mu.Lock()
					responses = append(responses, val)
					mu.Unlock()
					fmt.Printf("  [%s] 返回数据: %s (ts: %d)\n", n.ID, val.Value, val.Timestamp)
				} else {
					fmt.Printf("  [%s] 数据未找到\n", n.ID)
				}
			} else {
				fmt.Printf("  [%s] 离线，读取失败\n", n.ID)
			}
		}(node)
	}
	wg.Wait()

	if len(responses) >= R {
		// Read Repair logic: pick the latest version (LWW)
		latest := responses[0]
		for _, r := range responses {
			if r.Timestamp > latest.Timestamp {
				latest = r
			}
		}
		fmt.Printf("[Client] 读取成功 (响应数 %d >= R=%d)，选择最新值: %s\n", len(responses), R, latest.Value)
		return latest, true
	}

	fmt.Printf("[Client] 读取失败 (有效响应数 %d < R=%d)\n", len(responses), R)
	return Data{}, false
}

func main() {
	fmt.Println("=== 无主复制 (Leaderless Replication) 演示 ===")
	fmt.Println()
	fmt.Println("特点: 客户端直接发送请求到多个节点。使用法定人数 (Quorum) 保证一致性。")
	fmt.Println("参数: N (节点总数), W (写入确认数), R (读取确认数)")
	fmt.Println("关键规则: W + R > N (保证读写集合一定有交集，即读到最新值)")
	fmt.Println()

	N := 3
	W := 2
	R := 2
	cluster := NewCluster(N)

	// 1. 正常写入 (W=2)
	fmt.Println("--- 场景 1: 正常写入 ---")
	cluster.Write("config", "v1", W)
	fmt.Println()

	// 2. 一个节点挂掉时的写入 (仍满足 W=2)
	fmt.Println("--- 场景 2: Node-2 挂掉，尝试写入 ---")
	cluster.Nodes[2].Online = false
	cluster.Write("config", "v2", W)
	fmt.Println()

	// 3. 读取 (R=2)
	fmt.Println("--- 场景 3: 读取数据 (满足 R=2) ---")
	// 即使 Node-2 挂了且它是旧版本（或没有数据），只要我们读 R=2 个节点，
	// 就能从 Node-0 或 Node-1 中读到 v2。
	cluster.Read("config", R)
	fmt.Println()

	// 4. 更多节点挂掉，读取失败
	fmt.Println("--- 场景 4: Node-1 也挂掉，尝试读取 ---")
	cluster.Nodes[1].Online = false
	cluster.Read("config", R)

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("无主复制（如 Dynamo, Cassandra）提供了极高的可用性。")
	fmt.Println("权衡:")
	fmt.Println("- 优势: 允许 N-W 个节点故障仍可写入，允许 N-R 个节点故障仍可读取。")
	fmt.Println("- 挑战: 需要处理陈旧数据。常用技术包括读修复 (Read Repair) 和反熵 (Anti-entropy)。")
	fmt.Println("- 局限: 不支持复杂的事务。")
}

