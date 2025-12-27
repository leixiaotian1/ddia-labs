# DDIA 学习与实践指南

本指南旨在帮助你将《设计数据密集型应用》(DDIA) 书中的理论知识与本项目的 Go 语言实现代码相对照。通过“阅读理论 -> 运行 Demo -> 分析代码”的循环，可以更深入地理解分布式系统的设计权衡。

## 章节映射表

| DDIA 章节 | 核心知识点 | 对应项目 Lab | 重点关注文件 |
| :--- | :--- | :--- | :--- |
| **第 3 章：存储与检索** | 存储引擎、SSTables、LSM-Tree、B-Tree | [01-storage-engine](../labs/01-storage-engine/) | `lsm/main.go`, `btree/main.go` |
| | 哈希索引、B-Tree 索引 | [02-indexing](../labs/02-indexing/) | `hash-index/main.go` |
| **第 5 章：复制** | 主从复制、异步 vs 同步、复制延迟 | [04-replication](../labs/04-replication/) | `master-slave/main.go` |
| | 多主复制、冲突处理 (LWW) | [04-replication](../labs/04-replication/) | `multi-master/main.go` |
| | 无主复制、Quorum (W+R>N)、读修复 | [04-replication](../labs/04-replication/) | `leaderless/main.go` |
| | 因果关系、并发冲突、CRDT | [10-crdt](../labs/10-crdt/) | `g-counter/main.go` |
| **第 6 章：分区** | 键值范围分区、哈希分区、热点问题 | [05-partitioning](../labs/05-partitioning/) | `range-partition/main.go` |
| | 一致性哈希、动态再平衡 | [05-partitioning](../labs/05-partitioning/) | `consistent-hash/main.go` |
| **第 7 章：事务** | 隔离级别（脏读、不可重复读、幻读） | [03-isolation-levels](../labs/03-isolation-levels/) | `dirty-read/main.go` 等 |
| | 读已提交、可重复读、可串行化 | [03-isolation-levels](../labs/03-isolation-levels/) | `serializable/main.go` |
| **第 8 章：分布式系统的麻烦** | 时钟同步、网络延迟、不可靠网络 | [11-clock-sync](../labs/11-clock-sync/) | `logical-clock/main.go` |
| **第 9 章：一致性与共识** | 线性一致性 (Linearizability) | [07-consistency](../labs/07-consistency/) | `linearizability/main.go` |
| | 全序关系、因果一致性 | [07-consistency](../labs/07-consistency/) | `causal/main.go` |
| | 两阶段提交 (2PC)、原子提交 | [06-distributed-tx](../labs/06-distributed-tx/) | `two-phase-commit/main.go` |
| | Raft、Paxos、选主、共识达成 | [08-consensus](../labs/08-consensus/) | `raft/main.go` |
| **第 10 章：批处理** | MapReduce 编程模型 | [12-batch-processing](../labs/12-batch-processing/) | `mapreduce/main.go` |
| **第 11 章：流处理** | 事件流、滑动窗口聚合 | [13-stream-processing](../labs/13-stream-processing/) | `event-stream/main.go` |

## 建议学习路径

建议按照以下阶段进行学习：

### 第一阶段：单机基础 (Foundations)
重点理解数据是如何存储在磁盘上的，以及如何在高并发下保证正确性。
1. **隔离级别**：首先运行 `03-isolation-levels`。这是最直观的，你可以看到“破坏”隔离级别后会发生什么。
2. **存储结构**：阅读并运行 `01-storage-engine`。对比 LSM-tree (适合写) 和 B-tree (适合读) 的代码实现差异。

### 第二阶段：数据分布 (Distributed Data)
理解数据是如何在多台机器间流动的。
1. **复制**：运行 `04-replication`。重点观察 `leaderless` 模式下的 W+R>N 逻辑。
2. **分区**：运行 `05-partitioning`。手动尝试在一致性哈希代码中添加/删除节点，观察数据路由的变化。

### 第三阶段：一致性与共识 (Consistency & Consensus)
这是全书最难也是最精华的部分。
1. **一致性模型**：通过 `07-consistency` 理解什么是“线性化”，以及它与“最终一致性”在代码表现上的差异。
2. **共识算法**：运行 `08-consensus/raft`。仔细观察 Leader 挂掉后，集群是如何在 1-2 秒内自动选出新 Leader 的。

### 第四阶段：综合与扩展 (Advanced Topics)
1. **综合实践**：阅读 `14-simple-db`。看看一个简单的数据库是如何把存储、索引和锁事务结合在一起的。
2. **离线与流式计算**：通过 `12` 和 `13` 理解大数据处理的两种基本范式。

## 如何利用代码进行深度学习？

1. **修改参数**：例如在 `04-replication/leaderless` 中，尝试将 `W+R` 设置为小于或等于 `N`，看看是否还能保证读到最新值。
2. **注入故障**：在 `08-consensus/raft` 代码中，手动调用 `Stop()` 停掉两个节点，观察只有 1 个节点存活时系统是否还能正常工作。
3. **性能测试**：运行 `01-storage-engine/benchmark`，观察不同写入量下两类引擎的耗时曲线。

---
*愿你在分布式系统的海洋中航行顺利！*

