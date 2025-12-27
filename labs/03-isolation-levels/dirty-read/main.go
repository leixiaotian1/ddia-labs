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

// 模拟数据库
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

// Read Uncommitted 级别：允许读取未提交的数据
// 问题：脏读 - 事务A读取了事务B未提交的修改

func main() {
	fmt.Println("=== 脏读 (Dirty Read) 演示 ===")
	fmt.Println()
	fmt.Println("隔离级别: Read Uncommitted")
	fmt.Println("问题: 事务A读取了事务B未提交的修改")
	fmt.Println()

	db := NewDatabase()
	db.CreateAccount(1, 1000)

	var wg sync.WaitGroup

	// 事务A：转账操作（未提交）
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("[事务A] 开始转账：从账户1转出500")

		account := db.GetAccount(1)
		if account == nil {
			return
		}

		// 模拟事务处理时间
		time.Sleep(100 * time.Millisecond)

		// 修改余额（但未提交）
		account.mu.Lock()
		oldBalance := account.Balance
		account.Balance -= 500
		fmt.Printf("[事务A] 修改余额: %d -> %d (未提交)\n", oldBalance, account.Balance)
		account.mu.Unlock()

		// 模拟一些处理时间
		time.Sleep(200 * time.Millisecond)

		// 发生错误，回滚
		fmt.Println("[事务A] 发生错误，回滚事务")
		account.mu.Lock()
		account.Balance = oldBalance
		fmt.Printf("[事务A] 回滚后余额: %d\n", account.Balance)
		account.mu.Unlock()
	}()

	// 事务B：读取账户余额（在Read Uncommitted级别下可能读到未提交的数据）
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(150 * time.Millisecond) // 在事务A修改后、回滚前读取

		account := db.GetAccount(1)
		if account == nil {
			return
		}

		account.mu.RLock()
		balance := account.Balance
		account.mu.RUnlock()

		fmt.Printf("[事务B] 读取账户1余额: %d (脏读！读取了未提交的数据)\n", balance)
		fmt.Println("[事务B] 基于错误的数据做出决策...")
	}()

	wg.Wait()

	fmt.Println("\n=== 问题分析 ===")
	fmt.Println("在Read Uncommitted级别下：")
	fmt.Println("1. 事务B读取了事务A未提交的修改（余额=500）")
	fmt.Println("2. 事务A随后回滚，实际余额仍为1000")
	fmt.Println("3. 事务B基于错误的数据（500）做出了错误的决策")
	fmt.Println("\n解决方案：使用Read Committed或更高级别的隔离级别")
}
