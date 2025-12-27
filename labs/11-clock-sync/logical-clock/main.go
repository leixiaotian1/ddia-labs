package main

import (
	"fmt"
	"time"
)

// LamportClock 逻辑时钟
type LamportClock struct {
	counter int
}

func (lc *LamportClock) Tick() {
	lc.counter++
}

func (lc *LamportClock) Send() int {
	lc.counter++
	return lc.counter
}

func (lc *LamportClock) Receive(receivedCounter int) {
	if receivedCounter > lc.counter {
		lc.counter = receivedCounter
	}
	lc.counter++
}

func main() {
	fmt.Println("=== Lamport 逻辑时钟演示 ===")
	fmt.Println()
	fmt.Println("特点: 不依赖物理时钟，仅通过事件顺序保证因果一致。")
	fmt.Println()

	nodeA := &LamportClock{}
	nodeB := &LamportClock{}

	fmt.Println("1. 节点 A 发生本地事件")
	nodeA.Tick()
	fmt.Printf("   Node A Clock: %d\n", nodeA.counter)

	fmt.Println("\n2. 节点 A 发送消息给 B")
	msgTerm := nodeA.Send()
	fmt.Printf("   Node A Clock (after send): %d, Message Term: %d\n", nodeA.counter, msgTerm)

	fmt.Println("\n3. 节点 B 接收消息")
	fmt.Printf("   Node B Clock (before receive): %d\n", nodeB.counter)
	nodeB.Receive(msgTerm)
	fmt.Printf("   Node B Clock (after receive): %d\n", nodeB.counter)

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("1. Lamport 时钟保证了: 如果 a -> b，那么 C(a) < C(b)。")
	fmt.Println("2. 局限性: 不能通过 C(a) < C(b) 推断出 a -> b (可能并发)。这就引出了更强大的向量时钟。")
}

