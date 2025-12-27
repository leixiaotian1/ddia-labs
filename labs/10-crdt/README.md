# CRDT（无冲突复制数据类型）

## 概述

本lab实现了CRDT（Conflict-free Replicated Data Types），用于在无主复制系统中实现最终一致性。

## 对应DDIA章节

- 第5章：复制
- 第5.4节：无主复制

## CRDT类型

### G-Counter（增长计数器）
- **特点**: 只能增长的计数器，支持多节点并发更新
- **应用**: 点赞数、访问计数等

### LWW-Register（最后写入获胜寄存器）
- **特点**: 保留最后一次写入的值
- **应用**: 配置信息、用户状态等

## 运行方式

```bash
# 运行G-Counter示例
cd g-counter
go run main.go

# 运行LWW-Register示例
cd lww-register
go run main.go
```

## 关键概念

- **交换律**: 操作顺序不影响最终结果
- **结合律**: 操作可以任意组合
- **幂等性**: 重复操作不影响结果

