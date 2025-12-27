package main

import (
	"fmt"
)

// 哈希索引实现（简化版）
type HashIndex struct {
	buckets []*Bucket
	size    int
}

type Bucket struct {
	entries []Entry
}

type Entry struct {
	Key   int
	Value string
}

func NewHashIndex(size int) *HashIndex {
	buckets := make([]*Bucket, size)
	for i := range buckets {
		buckets[i] = &Bucket{entries: make([]Entry, 0)}
	}
	return &HashIndex{
		buckets: buckets,
		size:    size,
	}
}

// 哈希函数（简单取模）
func (hi *HashIndex) hash(key int) int {
	return key % hi.size
}

// 插入
func (hi *HashIndex) Put(key int, value string) {
	bucketIndex := hi.hash(key)
	bucket := hi.buckets[bucketIndex]

	// 检查key是否已存在
	for i, entry := range bucket.entries {
		if entry.Key == key {
			// 更新现有值
			bucket.entries[i].Value = value
			return
		}
	}

	// 添加新条目
	bucket.entries = append(bucket.entries, Entry{Key: key, Value: value})
}

// 查找
func (hi *HashIndex) Get(key int) (string, bool) {
	bucketIndex := hi.hash(key)
	bucket := hi.buckets[bucketIndex]

	// 在bucket中线性查找（处理哈希冲突）
	for _, entry := range bucket.entries {
		if entry.Key == key {
			return entry.Value, true
		}
	}

	return "", false
}

// 删除
func (hi *HashIndex) Delete(key int) bool {
	bucketIndex := hi.hash(key)
	bucket := hi.buckets[bucketIndex]

	for i, entry := range bucket.entries {
		if entry.Key == key {
			// 删除条目
			bucket.entries = append(bucket.entries[:i], bucket.entries[i+1:]...)
			return true
		}
	}

	return false
}

func main() {
	fmt.Println("=== 哈希索引演示 ===\n")
	fmt.Println("哈希索引特点：")
	fmt.Println("1. O(1)平均时间复杂度查找")
	fmt.Println("2. 不支持范围查询")
	fmt.Println("3. 数据无序存储")
	fmt.Println("4. 需要处理哈希冲突\n")

	index := NewHashIndex(10)

	// 插入数据
	fmt.Println("插入数据：")
	keys := []int{10, 20, 5, 15, 25, 30, 8, 12}
	values := []string{"val10", "val20", "val5", "val15", "val25", "val30", "val8", "val12"}

	for i, key := range keys {
		index.Put(key, values[i])
		fmt.Printf("  插入 key=%d, value=%s (hash=%d)\n", key, values[i], index.hash(key))
	}

	// 查找数据
	fmt.Println("\n查找数据：")
	testKeys := []int{10, 15, 25, 100}
	for _, key := range testKeys {
		value, found := index.Get(key)
		if found {
			fmt.Printf("  key=%d -> value=%s (查找次数: 1次哈希计算 + 线性查找)\n", key, value)
		} else {
			fmt.Printf("  key=%d -> 未找到\n", key)
		}
	}

	// 尝试范围查询（不支持）
	fmt.Println("\n范围查询 [10, 25]：")
	fmt.Println("  哈希索引不支持范围查询！")
	fmt.Println("  需要扫描所有bucket，效率低下")

	fmt.Println("\n=== 哈希索引权衡分析 ===")
	fmt.Println("优势：")
	fmt.Println("- 查找速度快，O(1)平均时间复杂度")
	fmt.Println("- 实现简单")
	fmt.Println("- 适合等值查询")
	fmt.Println("\n劣势：")
	fmt.Println("- 不支持范围查询")
	fmt.Println("- 数据无序，无法有序遍历")
	fmt.Println("- 需要处理哈希冲突")
	fmt.Println("- 哈希表大小固定或需要动态扩容")
	fmt.Println("\n适用场景：")
	fmt.Println("- 等值查询（WHERE key = ?）")
	fmt.Println("- 不需要范围查询的场景")
	fmt.Println("- 内存数据库的索引（如Redis）")
}

