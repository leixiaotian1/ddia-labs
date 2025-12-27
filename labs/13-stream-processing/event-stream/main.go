package main

import (
	"fmt"
	"time"
)

// Event 代表流中的一个事件
type Event struct {
	Timestamp time.Time
	Value     float64
}

func main() {
	fmt.Println("=== 事件流处理 (Stream Processing) 演示 ===")
	fmt.Println()
	fmt.Println("场景: 实时计算滑动窗口内的平均值。")
	fmt.Println()

	stream := make(chan Event)

	// 模拟流生成器
	go func() {
		data := []float64{10, 20, 30, 40, 50, 60}
		for _, d := range data {
			stream <- Event{Timestamp: time.Now(), Value: d}
			time.Sleep(100 * time.Millisecond)
		}
		close(stream)
	}()

	// 模拟流处理器 (Windowed Average)
	fmt.Println("--- 开始实时处理 ---")
	window := []float64{}
	windowSize := 3

	for event := range stream {
		window = append(window, event.Value)
		if len(window) > windowSize {
			window = window[1:]
		}

		sum := 0.0
		for _, v := range window {
			sum += v
		}
		avg := sum / float64(len(window))

		fmt.Printf("[Processor] 接收值: %.1f, 滑动窗口平均值: %.1f\n", event.Value, avg)
	}

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("1. 流处理是实时的、连续的，与离线的批处理形成对比。")
	fmt.Println("2. 滑动窗口是流处理中非常常见的聚合手段。")
	fmt.Println("3. 流处理系统（如 Flink, Kafka Streams）需要处理乱序、水位线 (Watermark) 和 Exactly-once 等复杂问题。")
}

