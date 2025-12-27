# 存储引擎对比

## 概述

本lab实现了两种主要的存储引擎结构，并对比它们的性能特征和适用场景：
- B-tree（B树）
- LSM-tree（Log-Structured Merge Tree）

## 对应DDIA章节

- 第3章：存储与检索
- 第3.1节：数据库的存储引擎
- 第3.2节：B-tree索引
- 第3.3节：LSM-tree

## 存储引擎对比

### B-tree
- **特点**: 原地更新，保持数据有序
- **优势**: 读取性能稳定，支持范围查询
- **劣势**: 写入可能产生随机I/O，需要维护索引

### LSM-tree
- **特点**: 追加写入，定期合并
- **优势**: 写入性能高（顺序写入），适合写多读少场景
- **劣势**: 读取可能需要查询多个层级，压缩可能影响写入

## 运行方式

```bash
# 运行B-tree示例
cd btree
go run main.go

# 运行LSM-tree示例
cd lsm
go run main.go

# 运行性能对比测试
cd benchmark
go run main.go
```

## 关键权衡

- **读取 vs 写入**: B-tree读取快，LSM-tree写入快
- **空间放大**: LSM-tree需要额外的存储空间（多版本数据）
- **写放大**: LSM-tree的压缩过程会产生写放大
- **适用场景**: 
  - B-tree: OLTP系统，需要快速随机读取
  - LSM-tree: 时序数据，日志系统，写密集型应用

