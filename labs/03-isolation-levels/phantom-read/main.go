package main

import (
	"fmt"
	"sync"
	"time"
)

// 模拟数据库中的订单表
type Order struct {
	ID     int
	UserID int
	Amount int
}

// 模拟数据库（Repeatable Read级别）
type Database struct {
	orders []*Order
	mu     sync.RWMutex
	nextID int
}

func NewDatabase() *Database {
	return &Database{
		orders: make([]*Order, 0),
		nextID: 1,
	}
}

func (db *Database) CreateOrder(userID int, amount int) *Order {
	db.mu.Lock()
	defer db.mu.Unlock()
	order := &Order{
		ID:     db.nextID,
		UserID: userID,
		Amount: amount,
	}
	db.orders = append(db.orders, order)
	db.nextID++
	return order
}

func (db *Database) CountOrdersByUser(userID int) int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	count := 0
	for _, order := range db.orders {
		if order.UserID == userID {
			count++
		}
	}
	return count
}

func (db *Database) GetOrdersByUser(userID int) []*Order {
	db.mu.RLock()
	defer db.mu.RUnlock()
	result := make([]*Order, 0)
	for _, order := range db.orders {
		if order.UserID == userID {
			result = append(result, order)
		}
	}
	return result
}

// Repeatable Read 级别：保证同一事务中多次读取同一行数据的一致性
// 问题：幻读 - 在同一事务中，范围查询的结果集发生变化

func main() {
	fmt.Println("=== 幻读 (Phantom Read) 演示 ===")
	fmt.Println()
	fmt.Println("隔离级别: Repeatable Read")
	fmt.Println("问题: 在同一事务中，范围查询的结果集发生变化")
	fmt.Println()

	db := NewDatabase()
	// 初始化：用户1已有2个订单
	db.CreateOrder(1, 100)
	db.CreateOrder(1, 200)

	var wg sync.WaitGroup

	// 事务A：统计用户1的订单数量和总金额
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("[事务A] 开始统计用户1的订单")

		// 第一次查询：统计订单数量
		count1 := db.CountOrdersByUser(1)
		fmt.Printf("[事务A] 第一次查询：用户1有 %d 个订单\n", count1)

		// 获取订单列表
		orders1 := db.GetOrdersByUser(1)
		sum1 := 0
		for _, order := range orders1 {
			sum1 += order.Amount
		}
		fmt.Printf("[事务A] 第一次查询：订单总金额 = %d\n", sum1)

		// 模拟一些处理时间
		time.Sleep(200 * time.Millisecond)

		// 第二次查询：再次统计订单数量
		count2 := db.CountOrdersByUser(1)
		fmt.Printf("[事务A] 第二次查询：用户1有 %d 个订单\n", count2)

		// 再次获取订单列表
		orders2 := db.GetOrdersByUser(1)
		sum2 := 0
		for _, order := range orders2 {
			sum2 += order.Amount
		}
		fmt.Printf("[事务A] 第二次查询：订单总金额 = %d\n", sum2)

		if count1 != count2 {
			fmt.Printf("[事务A] 警告：两次查询的订单数量不一致！(%d != %d)\n", count1, count2)
			fmt.Println("[事务A] 这是幻读问题！出现了新的订单（幻影行）")
		}
		if sum1 != sum2 {
			fmt.Printf("[事务A] 警告：两次查询的总金额不一致！(%d != %d)\n", sum1, sum2)
		}
	}()

	// 事务B：为用户1创建新订单（在事务A两次查询之间提交）
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond) // 在事务A第一次查询后执行

		fmt.Println("[事务B] 开始为用户1创建新订单：金额300")
		order := db.CreateOrder(1, 300)
		fmt.Printf("[事务B] 创建订单成功：ID=%d, 用户ID=%d, 金额=%d\n", order.ID, order.UserID, order.Amount)

		// 提交事务
		fmt.Println("[事务B] 提交事务")
	}()

	wg.Wait()

	fmt.Println("\n=== 问题分析 ===")
	fmt.Println("在Repeatable Read级别下：")
	fmt.Println("1. 事务A第一次查询：用户1有2个订单，总金额=300")
	fmt.Println("2. 事务B创建新订单并提交，用户1现在有3个订单")
	fmt.Println("3. 事务A第二次查询：用户1有3个订单，总金额=600（出现了新行）")
	fmt.Println("4. 虽然单行数据可重复读，但范围查询的结果集发生了变化（幻读）")
	fmt.Println("\n解决方案：使用Serializable隔离级别，通过范围锁防止幻读")
}
