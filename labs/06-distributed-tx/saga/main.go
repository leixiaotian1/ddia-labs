package main

import (
	"fmt"
	"time"
)

// Saga 代表一个长事务中的一系列本地事务及其对应的补偿操作
type Step struct {
	Name       string
	Execute    func() bool
	Compensate func()
}

func main() {
	fmt.Println("=== Saga 模式演示 (分布式事务补偿机制) ===")
	fmt.Println()
	fmt.Println("特点: 最终一致性。将长事务分解为一系列本地事务，失败时反向执行补偿。")
	fmt.Println()

	// 定义业务流程: 订机票 -> 订酒店 -> 支付
	steps := []Step{
		{
			Name: "订机票",
			Execute: func() bool {
				fmt.Println("[Step 1] 订机票成功")
				return true
			},
			Compensate: func() {
				fmt.Println("[Compensate 1] 退订机票")
			},
		},
		{
			Name: "订酒店",
			Execute: func() bool {
				fmt.Println("[Step 2] 订酒店失败 (模拟无房)")
				return false
			},
			Compensate: func() {
				fmt.Println("[Compensate 2] 取消酒店预订")
			},
		},
		{
			Name: "支付",
			Execute: func() bool {
				fmt.Println("[Step 3] 支付成功")
				return true
			},
			Compensate: func() {
				fmt.Println("[Compensate 3] 退款")
			},
		},
	}

	fmt.Println("--- 执行 Saga 流程 ---")
	history := []int{} // 记录执行成功的步骤索引

	success := true
	for i, step := range steps {
		if step.Execute() {
			history = append(history, i)
		} else {
			fmt.Printf("[System] 步骤 '%s' 失败，启动补偿流程...\n", step.Name)
			success = false
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !success {
		// 反向执行补偿操作
		for i := len(history) - 1; i >= 0; i-- {
			steps[history[i]].Compensate()
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println("[System] 补偿流程完成，系统恢复一致状态")
	} else {
		fmt.Println("[System] 事务执行成功")
	}

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("1. Saga 放弃了 ACID 中的隔离性 (I)，提高了系统性能和可用性。")
	fmt.Println("2. 适用于微服务架构下的长事务。")
	fmt.Println("3. 核心挑战是必须实现可靠的补偿操作。")
}

