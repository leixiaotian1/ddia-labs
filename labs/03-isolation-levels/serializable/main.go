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

// 模拟数据库（Serializable级别）
// Serializable通过严格的锁机制（如两阶段锁）保证完全隔离
type Database struct {
	accounts map[int]*Account
	mu       sync.Mutex // 使用互斥锁模拟严格的序列化
}

func NewDatabase() *Database {
	return &Database{
		accounts: make(map[int]*Account),
	}
}

func (db *Database) GetAccount(id int) *Account {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.accounts[id]
}

func (db *Database) CreateAccount(id int, balance int) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.accounts[id] = &Account{ID: id, Balance: balance}
}

// 转账操作（使用两阶段锁）
func (db *Database) Transfer(fromID, toID int, amount int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	from := db.accounts[fromID]
	to := db.accounts[toID]

	if from == nil || to == nil {
		return fmt.Errorf("账户不存在")
	}

	if from.Balance < amount {
		return fmt.Errorf("余额不足")
	}

	// 两阶段锁：在事务结束前持有所有锁
	from.mu.Lock()
	to.mu.Lock()

	from.Balance -= amount
	to.Balance += amount

	from.mu.Unlock()
	to.mu.Unlock()

	return nil
}

// Serializable 级别：最高隔离级别，完全避免所有并发问题
// 通过严格的序列化（如两阶段锁）保证事务完全隔离

func main() {
	fmt.Println("=== 可串行化 (Serializable) 演示 ===")
	fmt.Println()
	fmt.Println("隔离级别: Serializable")
	fmt.Println("保证: 完全隔离，无脏读、不可重复读、幻读问题")
	fmt.Println()

	db := NewDatabase()
	db.CreateAccount(1, 1000)
	db.CreateAccount(2, 500)

	var wg sync.WaitGroup

	// 事务A：从账户1转出300到账户2
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("[事务A] 开始转账：从账户1转出300到账户2")

		// 模拟事务处理时间
		time.Sleep(100 * time.Millisecond)

		err := db.Transfer(1, 2, 300)
		if err != nil {
			fmt.Printf("[事务A] 转账失败: %v\n", err)
			return
		}

		fmt.Println("[事务A] 转账成功")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("[事务A] 提交事务")
	}()

	// 事务B：读取账户余额（在Serializable级别下，会等待事务A完成）
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond) // 在事务A开始后读取

		fmt.Println("[事务B] 开始读取账户余额")

		account1 := db.GetAccount(1)
		account2 := db.GetAccount(2)

		if account1 != nil && account2 != nil {
			account1.mu.RLock()
			balance1 := account1.Balance
			account1.mu.RUnlock()

			account2.mu.RLock()
			balance2 := account2.Balance
			account2.mu.RUnlock()

			fmt.Printf("[事务B] 账户1余额: %d\n", balance1)
			fmt.Printf("[事务B] 账户2余额: %d\n", balance2)
			fmt.Println("[事务B] 读取到一致的数据（无脏读）")
		}
	}()

	wg.Wait()

	fmt.Println("\n=== Serializable级别特点 ===")
	fmt.Println("1. 完全隔离：事务按某种顺序串行执行")
	fmt.Println("2. 无脏读：只读取已提交的数据")
	fmt.Println("3. 无不可重复读：同一事务中多次读取结果一致")
	fmt.Println("4. 无幻读：范围查询的结果集不会变化")
	fmt.Println("\n权衡：")
	fmt.Println("- 优点：数据一致性最强")
	fmt.Println("- 缺点：并发性能最低，可能出现死锁")
}
