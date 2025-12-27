package main

import (
	"fmt"
)

// B-tree索引节点（简化版）
type BTreeIndexNode struct {
	keys     []int
	values   []string
	children []*BTreeIndexNode
	isLeaf   bool
}

// B-tree索引实现
type BTreeIndex struct {
	root *BTreeIndexNode
}

func NewBTreeIndex() *BTreeIndex {
	return &BTreeIndex{
		root: &BTreeIndexNode{
			keys:   make([]int, 0),
			values: make([]string, 0),
			isLeaf: true,
		},
	}
}

// 查找
func (bti *BTreeIndex) Get(key int) (string, bool) {
	return bti.search(bti.root, key)
}

func (bti *BTreeIndex) search(node *BTreeIndexNode, key int) (string, bool) {
	if node == nil {
		return "", false
	}

	// 在节点中查找key的位置
	i := 0
	for i < len(node.keys) && key > node.keys[i] {
		i++
	}

	// 如果找到key
	if i < len(node.keys) && key == node.keys[i] {
		if node.isLeaf {
			return node.values[i], true
		}
	}

	// 如果是叶子节点但没找到，返回false
	if node.isLeaf {
		return "", false
	}

	// 递归查找子节点
	return bti.search(node.children[i], key)
}

// 插入
func (bti *BTreeIndex) Put(key int, value string) {
	bti.insert(bti.root, key, value)
}

func (bti *BTreeIndex) insert(node *BTreeIndexNode, key int, value string) {
	if node.isLeaf {
		// 在叶子节点中插入
		pos := 0
		for pos < len(node.keys) && key > node.keys[pos] {
			pos++
		}

		// 如果key已存在，更新value
		if pos < len(node.keys) && node.keys[pos] == key {
			node.values[pos] = value
			return
		}

		// 插入新key-value
		node.keys = append(node.keys[:pos], append([]int{key}, node.keys[pos:]...)...)
		node.values = append(node.values[:pos], append([]string{value}, node.values[pos:]...)...)
		return
	}

	// 非叶子节点，找到合适的子节点
	pos := 0
	for pos < len(node.keys) && key > node.keys[pos] {
		pos++
	}
	bti.insert(node.children[pos], key, value)
}

// 范围查询
func (bti *BTreeIndex) RangeScan(start, end int) []string {
	var result []string
	bti.rangeScan(bti.root, start, end, &result)
	return result
}

func (bti *BTreeIndex) rangeScan(node *BTreeIndexNode, start, end int, result *[]string) {
	if node == nil {
		return
	}

	if node.isLeaf {
		// 在叶子节点中查找范围内的key
		for i, key := range node.keys {
			if key >= start && key <= end {
				*result = append(*result, node.values[i])
			}
		}
		return
	}

	// 非叶子节点，递归查找
	for i, key := range node.keys {
		if key >= start {
			bti.rangeScan(node.children[i], start, end, result)
		}
		if key > end {
			break
		}
	}
	// 检查最后一个子节点
	if len(node.keys) > 0 && node.keys[len(node.keys)-1] < end {
		bti.rangeScan(node.children[len(node.children)-1], start, end, result)
	}
}

// 有序遍历
func (bti *BTreeIndex) InOrderTraversal() []string {
	var result []string
	bti.inOrder(bti.root, &result)
	return result
}

func (bti *BTreeIndex) inOrder(node *BTreeIndexNode, result *[]string) {
	if node == nil {
		return
	}

	if node.isLeaf {
		*result = append(*result, node.values...)
		return
	}

	// 遍历所有子节点和键
	for i := 0; i < len(node.children); i++ {
		bti.inOrder(node.children[i], result)
		if i < len(node.keys) {
			// 非叶子节点的值通常不存储实际数据
		}
	}
}

func main() {
	fmt.Println("=== B-tree索引演示 ===\n")
	fmt.Println("B-tree索引特点：")
	fmt.Println("1. O(log n)查找时间复杂度")
	fmt.Println("2. 支持范围查询")
	fmt.Println("3. 数据有序存储")
	fmt.Println("4. 支持有序遍历\n")

	index := NewBTreeIndex()

	// 插入数据
	fmt.Println("插入数据：")
	keys := []int{10, 20, 5, 15, 25, 30, 8, 12}
	values := []string{"val10", "val20", "val5", "val15", "val25", "val30", "val8", "val12"}

	for i, key := range keys {
		index.Put(key, values[i])
		fmt.Printf("  插入 key=%d, value=%s\n", key, values[i])
	}

	// 查找数据
	fmt.Println("\n查找数据：")
	testKeys := []int{10, 15, 25, 100}
	for _, key := range testKeys {
		value, found := index.Get(key)
		if found {
			fmt.Printf("  key=%d -> value=%s (O(log n)查找)\n", key, value)
		} else {
			fmt.Printf("  key=%d -> 未找到\n", key)
		}
	}

	// 范围查询
	fmt.Println("\n范围查询 [10, 25]：")
	results := index.RangeScan(10, 25)
	for _, value := range results {
		fmt.Printf("  %s\n", value)
	}
	fmt.Println("  B-tree索引支持高效的范围查询！")

	fmt.Println("\n=== B-tree索引权衡分析 ===")
	fmt.Println("优势：")
	fmt.Println("- 支持范围查询（WHERE key BETWEEN ? AND ?）")
	fmt.Println("- 数据有序，支持ORDER BY")
	fmt.Println("- 查找性能稳定，O(log n)")
	fmt.Println("- 支持前缀匹配")
	fmt.Println("\n劣势：")
	fmt.Println("- 查找速度略慢于哈希索引（O(log n) vs O(1)）")
	fmt.Println("- 需要维护树结构")
	fmt.Println("- 写入可能产生随机I/O")
	fmt.Println("\n适用场景：")
	fmt.Println("- 需要范围查询的场景")
	fmt.Println("- 需要有序遍历的场景")
	fmt.Println("- 大多数关系型数据库的默认索引（MySQL, PostgreSQL）")
	fmt.Println("- 需要支持ORDER BY的查询")
}

