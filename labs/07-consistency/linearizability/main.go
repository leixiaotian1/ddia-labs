package main

import (
	"fmt"
	"sync"
	"time"
)

// LinearizableStorage simulates a storage that should be linearizable
type LinearizableStorage struct {
	value string
	mu    sync.Mutex
}

func (s *LinearizableStorage) Write(val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.value = val
}

func (s *LinearizableStorage) Read() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.value
}

func main() {
	fmt.Println("=== 线性一致性 (Linearizability) 演示 ===")
	fmt.Println()
	fmt.Println("定义: 任何操作在调用和返回之间的某个时刻原子地生效。")
	fmt.Println("表现: 一旦某个读取返回了新值，所有后续读取必须返回该新值或更新的值。")
	fmt.Println()

	storage := &LinearizableStorage{value: "Initial"}

	var wg sync.WaitGroup

	// 1. 模拟写入者
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("[Writer] 开始写入 'Updated'...")
		time.Sleep(100 * time.Millisecond) // 模拟处理
		storage.Write("Updated")
		fmt.Println("[Writer] 写入完成")
	}()

	// 2. 模拟多个并发读取者
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(50+id*20) * time.Millisecond)
			val := storage.Read()
			fmt.Printf("[Reader-%d] 读取结果: %s\n", id, val)
		}(i)
	}

	wg.Wait()

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("在真正的线性化系统中，读取结果必须符合全序关系。")
	fmt.Println("注意: 线性化不同于可串行化 (Serializability)。线性化针对的是单个对象的读写一致性，而可串行化针对的是跨事务的隔离性。")
}

