package main

import (
	"fmt"
)

// B-tree节点
type BTreeNode struct {
	keys     []int
	values   []string
	children []*BTreeNode
	isLeaf   bool
}

// B-tree实现（简化版，仅用于演示）
type BTree struct {
	root  *BTreeNode
	order int // B-tree的阶数（每个节点最多order-1个键）
}

func NewBTree(order int) *BTree {
	return &BTree{
		root: &BTreeNode{
			keys:   make([]int, 0),
			values: make([]string, 0),
			isLeaf: true,
		},
		order: order,
	}
}

// 查找操作
func (bt *BTree) Get(key int) (string, bool) {
	return bt.search(bt.root, key)
}

func (bt *BTree) search(node *BTreeNode, key int) (string, bool) {
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
		// 非叶子节点，继续向下查找
	}

	// 如果是叶子节点但没找到，返回false
	if node.isLeaf {
		return "", false
	}

	// 递归查找子节点
	return bt.search(node.children[i], key)
}

// 插入操作（简化版，仅支持叶子节点插入）
func (bt *BTree) Put(key int, value string) {
	bt.insert(bt.root, key, value)
}

func (bt *BTree) insert(node *BTreeNode, key int, value string) {
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
	bt.insert(node.children[pos], key, value)
}

// 范围查询
func (bt *BTree) RangeScan(start, end int) []string {
	var result []string
	bt.rangeScan(bt.root, start, end, &result)
	return result
}

func (bt *BTree) rangeScan(node *BTreeNode, start, end int, result *[]string) {
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
			bt.rangeScan(node.children[i], start, end, result)
		}
		if key > end {
			break
		}
	}
	// 检查最后一个子节点
	if len(node.keys) > 0 && node.keys[len(node.keys)-1] < end {
		bt.rangeScan(node.children[len(node.children)-1], start, end, result)
	}
}

func main() {
	fmt.Println("=== B-tree 存储引擎演示 ===")
	fmt.Println()
	fmt.Println("B-tree特点：")
	fmt.Println("1. 数据按key有序存储")
	fmt.Println("2. 支持高效的随机读取")
	fmt.Println("3. 支持范围查询")
	fmt.Println("4. 写入时可能需要更新多个节点（随机I/O）")
	fmt.Println()

	bt := NewBTree(4)

	// 插入数据
	fmt.Println("插入数据：")
	keys := []int{10, 20, 5, 15, 25, 30, 8, 12}
	values := []string{"val10", "val20", "val5", "val15", "val25", "val30", "val8", "val12"}

	for i, key := range keys {
		bt.Put(key, values[i])
		fmt.Printf("  插入 key=%d, value=%s\n", key, values[i])
	}

	// 查找数据
	fmt.Println("\n查找数据：")
	testKeys := []int{10, 15, 100}
	for _, key := range testKeys {
		value, found := bt.Get(key)
		if found {
			fmt.Printf("  key=%d -> value=%s\n", key, value)
		} else {
			fmt.Printf("  key=%d -> 未找到\n", key)
		}
	}

	// 范围查询
	fmt.Println("\n范围查询 [10, 25]：")
	results := bt.RangeScan(10, 25)
	for _, value := range results {
		fmt.Printf("  %s\n", value)
	}

	fmt.Println("\n=== B-tree 权衡分析 ===")
	fmt.Println("优势：")
	fmt.Println("- 读取性能稳定，O(log n)时间复杂度")
	fmt.Println("- 支持高效的范围查询")
	fmt.Println("- 数据有序，便于遍历")
	fmt.Println("\n劣势：")
	fmt.Println("- 写入可能产生随机I/O（需要更新多个节点）")
	fmt.Println("- 需要维护索引结构，占用额外空间")
	fmt.Println("- 写入性能不如LSM-tree")
}
