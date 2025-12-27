# 事务隔离级别测试

## 概述

本lab演示了SQL标准定义的四个事务隔离级别，以及每个级别可能出现的并发问题。通过"破坏性"测试，直观地展示不同隔离级别下可能出现的问题。

## 对应DDIA章节

- 第7章：事务
- 第7.2节：弱隔离级别

## 隔离级别

### 1. Read Uncommitted (读未提交)
- **问题**: 脏读（Dirty Read）
- **演示**: `dirty-read/`

### 2. Read Committed (读已提交)
- **问题**: 不可重复读（Non-repeatable Read）
- **演示**: `non-repeatable/`

### 3. Repeatable Read (可重复读)
- **问题**: 幻读（Phantom Read）
- **演示**: `phantom-read/`

### 4. Serializable (可串行化)
- **保证**: 完全隔离，无并发问题
- **演示**: `serializable/`

## 运行方式

```bash
# 运行所有隔离级别测试
go run main.go

# 运行特定隔离级别的测试
cd dirty-read
go run main.go
```

## 预期结果

每个子目录的demo会展示：
1. 在低隔离级别下，如何重现并发问题
2. 问题发生的时序图
3. 如何通过提高隔离级别解决问题

## 关键权衡

- **隔离级别 vs 性能**: 隔离级别越高，并发性能越低
- **隔离级别 vs 正确性**: 隔离级别越低，可能出现的数据不一致问题越多

