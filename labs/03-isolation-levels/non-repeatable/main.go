package main

import (
	"fmt"
	"sync"
	"time"
)

// 模拟数据库中的账户表
type Account struct {
	ID      int
	Balance int
	mu      sync.RWMutex
}

// 模拟数据库（Read Committed级别）
type Database struct {
	accounts map[int]*Account
	mu       sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		accounts: make(map[int]*Account),
	}
}

func (db *Database) GetAccount(id int) *Account {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.accounts[id]
}

func (db *Database) CreateAccount(id int, balance int) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.accounts[id] = &Account{ID: id, Balance: balance}
}

// Read Committed 级别：只读取已提交的数据
// 问题：不可重复读 - 在同一事务中，两次读取同一数据得到不同的结果

func main() {
	fmt.Println("=== 不可重复读 (Non-repeatable Read) 演示 ===")
	fmt.Println()
	fmt.Println("隔离级别: Read Committed")
	fmt.Println("问题: 在同一事务中，两次读取同一数据得到不同的结果")
	fmt.Println()

	db := NewDatabase()
	db.CreateAccount(1, 1000)

	var wg sync.WaitGroup

	// 事务A：统计操作（需要读取两次）
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("[事务A] 开始统计操作：计算账户1的余额总和")

		account := db.GetAccount(1)
		if account == nil {
			return
		}

		// 第一次读取
		account.mu.RLock()
		balance1 := account.Balance
		account.mu.RUnlock()
		fmt.Printf("[事务A] 第一次读取账户1余额: %d\n", balance1)

		// 模拟一些处理时间
		time.Sleep(200 * time.Millisecond)

		// 第二次读取（在Read Committed级别下，可能读到不同的值）
		account.mu.RLock()
		balance2 := account.Balance
		account.mu.RUnlock()
		fmt.Printf("[事务A] 第二次读取账户1余额: %d\n", balance2)

		if balance1 != balance2 {
			fmt.Printf("[事务A] 警告：两次读取结果不一致！(%d != %d)\n", balance1, balance2)
			fmt.Println("[事务A] 这是不可重复读问题！")
		}
	}()

	// 事务B：更新账户余额（在事务A两次读取之间提交）
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond) // 在事务A第一次读取后执行

		account := db.GetAccount(1)
		if account == nil {
			return
		}

		fmt.Println("[事务B] 开始更新账户1余额：增加200")
		account.mu.Lock()
		oldBalance := account.Balance
		account.Balance += 200
		fmt.Printf("[事务B] 更新余额: %d -> %d\n", oldBalance, account.Balance)
		account.mu.Unlock()

		// 提交事务
		fmt.Println("[事务B] 提交事务")
	}()

	wg.Wait()

	fmt.Println("\n=== 问题分析 ===")
	fmt.Println("在Read Committed级别下：")
	fmt.Println("1. 事务A第一次读取账户1余额 = 1000")
	fmt.Println("2. 事务B更新并提交，余额变为 1200")
	fmt.Println("3. 事务A第二次读取账户1余额 = 1200（与第一次不同）")
	fmt.Println("4. 事务A在同一事务中看到了不同的数据，导致统计结果不一致")
	fmt.Println("\n解决方案：使用Repeatable Read或更高级别的隔离级别")
}
