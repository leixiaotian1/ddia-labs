package main

import (
	"fmt"
	"sort"
)

// Data represents a record
type Data struct {
	Key   string
	Value string
}

// Partition represents a physical storage unit
type Partition struct {
	ID    int
	Range [2]string // [Min, Max)
	Data  []Data
}

// RangePartitioner manages partitions by key ranges
type RangePartitioner struct {
	Partitions []*Partition
}

func NewRangePartitioner(boundaries []string) *RangePartitioner {
	rp := &RangePartitioner{}
	
	// Ensure boundaries are sorted
	sort.Strings(boundaries)
	
	last := ""
	for i, b := range boundaries {
		rp.Partitions = append(rp.Partitions, &Partition{
			ID:    i,
			Range: [2]string{last, b},
			Data:  []Data{},
		})
		last = b
	}
	// Add final partition
	rp.Partitions = append(rp.Partitions, &Partition{
		ID:    len(boundaries),
		Range: [2]string{last, "ZZZZZZ"}, // Simplified max
		Data:  []Data{},
	})
	
	return rp
}

func (rp *RangePartitioner) GetPartitionID(key string) int {
	for i, p := range rp.Partitions {
		if key >= p.Range[0] && (p.Range[1] == "" || key < p.Range[1]) {
			return i
		}
	}
	return len(rp.Partitions) - 1
}

func (rp *RangePartitioner) Write(key, value string) {
	id := rp.GetPartitionID(key)
	rp.Partitions[id].Data = append(rp.Partitions[id].Data, Data{Key: key, Value: value})
	fmt.Printf("[Write] Key='%s' -> Partition-%d (范围: [%s, %s))\n", 
		key, id, rp.Partitions[id].Range[0], rp.Partitions[id].Range[1])
}

func main() {
	fmt.Println("=== 范围分区 (Range Partitioning) 演示 ===")
	fmt.Println()
	fmt.Println("特点: 根据 Key 的字典序范围来划分数据。支持高效的范围查询。")
	fmt.Println("挑战: 容易产生热点（Hotspots）。例如，如果 Key 是时间戳。")
	fmt.Println()

	// 定义分区边界：[, D), [D, M), [M, T), [T, )
	rp := NewRangePartitioner([]string{"D", "M", "T"})

	fmt.Println("--- 数据写入 ---")
	keys := []string{"Apple", "Banana", "Cat", "Dog", "Egg", "Monkey", "Orange", "Tiger", "Zebra"}
	for _, k := range keys {
		rp.Write(k, "val_"+k)
	}

	fmt.Println()
	fmt.Println("--- 范围查询演示 ---")
	fmt.Println("查询 [A, E) 范围内的数据:")
	// 在范围分区中，我们只需要查询 Partition-0 和 Partition-1
	start, end := "A", "E"
	for _, p := range rp.Partitions {
		// 检查分区范围是否与查询范围有交集
		if !(p.Range[1] <= start || p.Range[0] >= end) {
			fmt.Printf("  检查 Partition-%d (%s to %s):\n", p.ID, p.Range[0], p.Range[1])
			for _, d := range p.Data {
				if d.Key >= start && d.Key < end {
					fmt.Printf("    找到: %s\n", d.Key)
				}
			}
		}
	}

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("优势: 对于范围查询（如 WHERE key > 'A' AND key < 'E'）非常友好，只需要扫描特定分区。")
	fmt.Println("劣势: 如果某些范围的数据特别多，会造成负载不均。")
}

