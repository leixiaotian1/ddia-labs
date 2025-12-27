package main

import (
	"crypto/md5"
	"fmt"
	"math/big"
)

// Partition represents a physical storage unit
type Partition struct {
	ID   int
	Data map[string]string
}

// HashPartitioner manages partitions by hashing keys
type HashPartitioner struct {
	Partitions []*Partition
	Count      int
}

func NewHashPartitioner(count int) *HashPartitioner {
	p := make([]*Partition, count)
	for i := 0; i < count; i++ {
		p[i] = &Partition{ID: i, Data: make(map[string]string)}
	}
	return &HashPartitioner{Partitions: p, Count: count}
}

// hashKey returns an integer representation of the key's hash
func hashKey(key string) *big.Int {
	h := md5.New()
	h.Write([]byte(key))
	return new(big.Int).SetBytes(h.Sum(nil))
}

func (hp *HashPartitioner) GetPartitionID(key string) int {
	h := hashKey(key)
	// partition = hash % count
	return int(new(big.Int).Mod(h, big.NewInt(int64(hp.Count))).Int64())
}

func (hp *HashPartitioner) Write(key, value string) {
	id := hp.GetPartitionID(key)
	hp.Partitions[id].Data[key] = value
	fmt.Printf("[Write] Key='%s' -> HashPart-%d\n", key, id)
}

func main() {
	fmt.Println("=== 哈希分区 (Hash Partitioning) 演示 ===")
	fmt.Println()
	fmt.Println("特点: 通过哈希函数将 Key 均匀分散到不同分区。有效防止热点。")
	fmt.Println("挑战: 丧失了范围查询能力。分区扩容时数据迁移量大（本例 hash % N）。")
	fmt.Println()

	hp := NewHashPartitioner(4)

	fmt.Println("--- 数据写入 ---")
	keys := []string{"Apple", "Banana", "Cat", "Dog", "Egg", "Monkey", "Orange", "Tiger", "Zebra"}
	for _, k := range keys {
		hp.Write(k, "val_"+k)
	}

	fmt.Println()
	fmt.Println("--- 分区分布统计 ---")
	for _, p := range hp.Partitions {
		fmt.Printf("Partition-%d: %d 个记录\n", p.ID, len(p.Data))
		for k := range p.Data {
			fmt.Printf("  - %s\n", k)
		}
	}

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("优势: 即使 Key 序列是有序的，哈希后也会均匀分布，负载均衡好。")
	fmt.Println("劣势: 不支持高效范围查询。改变分区数量 N 时，几乎所有数据都需要重新路由（迁移）。")
}

