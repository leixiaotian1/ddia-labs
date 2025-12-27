package main

import (
	"fmt"
	"sort"
)

// LSM-tree的MemTable（内存表）
type MemTable struct {
	data map[int]string
}

func NewMemTable() *MemTable {
	return &MemTable{
		data: make(map[int]string),
	}
}

func (mt *MemTable) Put(key int, value string) {
	mt.data[key] = value
}

func (mt *MemTable) Get(key int) (string, bool) {
	value, ok := mt.data[key]
	return value, ok
}

func (mt *MemTable) Size() int {
	return len(mt.data)
}

// SSTable（Sorted String Table）- 磁盘上的有序表
type SSTable struct {
	entries []Entry
}

type Entry struct {
	Key   int
	Value string
}

func NewSSTable(entries []Entry) *SSTable {
	// 确保有序
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return &SSTable{entries: entries}
}

func (sst *SSTable) Get(key int) (string, bool) {
	// 二分查找
	left, right := 0, len(sst.entries)-1
	for left <= right {
		mid := (left + right) / 2
		if sst.entries[mid].Key == key {
			return sst.entries[mid].Value, true
		} else if sst.entries[mid].Key < key {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return "", false
}

// LSM-tree实现（简化版）
type LSMTree struct {
	memTable  *MemTable
	sstables  []*SSTable // 多个层级的SSTable
	threshold int        // MemTable刷新阈值
}

func NewLSMTree(threshold int) *LSMTree {
	return &LSMTree{
		memTable:  NewMemTable(),
		sstables:  make([]*SSTable, 0),
		threshold: threshold,
	}
}

// 写入操作（追加写入）
func (lsm *LSMTree) Put(key int, value string) {
	lsm.memTable.Put(key, value)

	// 如果MemTable达到阈值，刷新到磁盘（这里简化为创建SSTable）
	if lsm.memTable.Size() >= lsm.threshold {
		lsm.flushMemTable()
	}
}

func (lsm *LSMTree) flushMemTable() {
	fmt.Println("  [LSM] MemTable达到阈值，刷新到SSTable...")

	// 将MemTable转换为SSTable
	entries := make([]Entry, 0, lsm.memTable.Size())
	for key, value := range lsm.memTable.data {
		entries = append(entries, Entry{Key: key, Value: value})
	}

	sstable := NewSSTable(entries)
	lsm.sstables = append(lsm.sstables, sstable)

	// 清空MemTable
	lsm.memTable = NewMemTable()

	fmt.Printf("  [LSM] 创建SSTable，包含 %d 个条目\n", len(entries))
}

// 读取操作（需要检查MemTable和所有SSTable）
func (lsm *LSMTree) Get(key int) (string, bool) {
	// 先检查MemTable
	if value, ok := lsm.memTable.Get(key); ok {
		return value, true
	}

	// 从最新的SSTable开始查找（LSM-tree通常从新到旧查找）
	for i := len(lsm.sstables) - 1; i >= 0; i-- {
		if value, ok := lsm.sstables[i].Get(key); ok {
			return value, true
		}
	}

	return "", false
}

func main() {
	fmt.Println("=== LSM-tree 存储引擎演示 ===")
	fmt.Println()
	fmt.Println("LSM-tree特点：")
	fmt.Println("1. 写入采用追加方式（顺序I/O）")
	fmt.Println("2. 数据在内存中积累，定期刷新到磁盘")
	fmt.Println("3. 读取需要检查多个层级（MemTable + SSTables）")
	fmt.Println("4. 定期合并（Compaction）减少读取开销")
	fmt.Println()

	lsm := NewLSMTree(5) // 阈值设为5

	// 插入数据
	fmt.Println("插入数据：")
	keys := []int{10, 20, 5, 15, 25, 30, 8, 12, 18, 22}
	values := []string{"val10", "val20", "val5", "val15", "val25", "val30", "val8", "val12", "val18", "val22"}

	for i, key := range keys {
		lsm.Put(key, values[i])
		fmt.Printf("  插入 key=%d, value=%s\n", key, values[i])
	}

	// 查找数据
	fmt.Println("\n查找数据：")
	testKeys := []int{10, 15, 25, 100}
	for _, key := range testKeys {
		value, found := lsm.Get(key)
		if found {
			fmt.Printf("  key=%d -> value=%s\n", key, value)
		} else {
			fmt.Printf("  key=%d -> 未找到\n", key)
		}
	}

	fmt.Println("\n=== LSM-tree 权衡分析 ===")
	fmt.Println("优势：")
	fmt.Println("- 写入性能高（顺序I/O，追加写入）")
	fmt.Println("- 适合写密集型场景")
	fmt.Println("- 写入吞吐量高")
	fmt.Println("\n劣势：")
	fmt.Println("- 读取可能需要查询多个SSTable（读放大）")
	fmt.Println("- 需要定期合并（Compaction），产生写放大")
	fmt.Println("- 空间放大（多版本数据同时存在）")
	fmt.Println("\n适用场景：")
	fmt.Println("- 时序数据库（InfluxDB, TimescaleDB）")
	fmt.Println("- 日志存储系统")
	fmt.Println("- 写多读少的应用")
}
