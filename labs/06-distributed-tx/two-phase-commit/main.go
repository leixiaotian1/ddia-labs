package main

import (
	"fmt"
	"sync"
	"time"
)

// 2PC 参与者状态
const (
	Ready   = "Ready"
	Aborted = "Aborted"
	Committed = "Committed"
)

type Participant struct {
	ID     string
	State  string
	mu     sync.Mutex
	Online bool
}

func (p *Participant) Prepare() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.Online {
		return false
	}
	fmt.Printf("[Participant-%s] 准备就绪 (Ready)\n", p.ID)
	p.State = Ready
	return true
}

func (p *Participant) Commit() {
	p.mu.Lock()
	defer p.mu.Unlock()
	fmt.Printf("[Participant-%s] 提交成功 (Committed)\n", p.ID)
	p.State = Committed
}

func (p *Participant) Rollback() {
	p.mu.Lock()
	defer p.mu.Unlock()
	fmt.Printf("[Participant-%s] 回滚 (Rollback)\n", p.ID)
	p.State = Aborted
}

func main() {
	fmt.Println("=== 两阶段提交 (Two-Phase Commit, 2PC) 演示 ===")
	fmt.Println()
	fmt.Println("特点: 强原子性，由协调者 (Coordinator) 驱动参与者。")
	fmt.Println("注意: 2PC 是阻塞协议，存在单点故障风险。")
	fmt.Println()

	participants := []*Participant{
		{ID: "Bank-A", Online: true},
		{ID: "Bank-B", Online: true},
	}

	// 1. 成功场景
	fmt.Println("--- 场景 1: 成功提交 ---")
	allReady := true
	for _, p := range participants {
		if !p.Prepare() {
			allReady = false
			break
		}
	}

	if allReady {
		fmt.Println("[Coordinator] 所有参与者已 Ready，发送 Commit 指令")
		for _, p := range participants {
			p.Commit()
		}
	} else {
		fmt.Println("[Coordinator] 有参与者失败，发送 Rollback 指令")
		for _, p := range participants {
			p.Rollback()
		}
	}

	fmt.Println()

	// 2. 失败场景
	fmt.Println("--- 场景 2: 参与者挂掉导致回滚 ---")
	participants[1].Online = false // 模拟 Bank-B 故障
	
	allReady = true
	for _, p := range participants {
		if !p.Prepare() {
			allReady = false
			fmt.Printf("[Coordinator] 发现 %s 无法 Ready\n", p.ID)
			break
		}
	}

	if allReady {
		fmt.Println("[Coordinator] 发送 Commit")
		for _, p := range participants {
			p.Commit()
		}
	} else {
		fmt.Println("[Coordinator] 有参与者失败，执行原子回滚")
		for _, p := range participants {
			p.Rollback()
		}
	}

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("1. 2PC 保证了分布式系统中的原子性（要么全做，要么全不做）。")
	fmt.Println("2. 核心缺陷是锁定资源时间长，且如果协调者崩溃，参与者可能无限期阻塞。")
}

