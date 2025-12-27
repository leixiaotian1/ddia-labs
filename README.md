# DDIA Labs

通过Go语言实现的实践demo，深入理解《设计数据密集型应用》(Designing Data-Intensive Applications) 中的核心概念和分布式系统的权衡决策。

## 项目目标

本项目旨在通过可运行的代码示例，帮助理解DDIA书中提到的：
- 数据存储与检索机制
- 数据复制与分区策略
- 事务与隔离级别
- 分布式系统的一致性模型
- 共识算法
- 批处理与流处理

## 项目结构

```
ddia-labs/
├── docs/          # 文档目录
├── pkg/           # 共享库
└── labs/          # 实验demo目录
```

## Labs列表

### Phase 1 - 基础存储与事务（核心）
- [01-storage-engine](labs/01-storage-engine/) - 存储引擎对比（B-tree vs LSM-tree）
- [02-indexing](labs/02-indexing/) - 索引结构（哈希索引、B-tree索引）
- [03-isolation-levels](labs/03-isolation-levels/) - 事务隔离级别"破坏性"测试 ⭐

### Phase 2 - 分布式基础
- [04-replication](labs/04-replication/) - 数据复制（主从、多主、无主）
- [05-partitioning](labs/05-partitioning/) - 数据分区策略
- [07-consistency](labs/07-consistency/) - 一致性模型对比

### Phase 3 - 高级主题
- [06-distributed-tx](labs/06-distributed-tx/) - 分布式事务
- [08-consensus](labs/08-consensus/) - 共识算法（Raft、Paxos）
- [14-simple-db](labs/14-simple-db/) - 简单数据库综合实现

### Phase 4 - 扩展主题
- [09-vector-clock](labs/09-vector-clock/) - 向量时钟
- [10-crdt](labs/10-crdt/) - CRDT（无冲突复制数据类型）
- [11-clock-sync](labs/11-clock-sync/) - 时钟同步问题
- [12-batch-processing](labs/12-batch-processing/) - 批处理（MapReduce）
- [13-stream-processing](labs/13-stream-processing/) - 流处理

## 快速开始

### 前置要求
- Go 1.21 或更高版本

### 运行示例

每个lab都是独立的，可以单独运行：

```bash
# 进入某个lab目录
cd labs/03-isolation-levels

# 运行示例
go run main.go

# 运行测试
go test ./...
```

## 学习路径建议

1. **入门**: 从 `03-isolation-levels` 开始，理解事务隔离级别
2. **存储基础**: 学习 `01-storage-engine` 和 `02-indexing`
3. **分布式基础**: 探索 `04-replication` 和 `05-partitioning`
4. **高级主题**: 深入 `08-consensus` 和 `14-simple-db`

## 贡献

欢迎提交Issue和Pull Request！

## 参考

- [Designing Data-Intensive Applications](https://dataintensive.net/) by Martin Kleppmann
- [DDIA 中文翻译](https://github.com/Vonng/ddia)
