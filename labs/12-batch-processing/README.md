# 批处理

## 概述

本lab实现了MapReduce的简化版本，演示批处理系统的核心概念。

## 对应DDIA章节

- 第10章：批处理
- 第10.1节：使用Unix工具的批处理
- 第10.2节：MapReduce与分布式文件系统

## MapReduce

MapReduce是一种编程模型，用于处理和生成大数据集。它将计算分为Map和Reduce两个阶段。

## 运行方式

```bash
cd mapreduce
go run main.go
```

## 关键概念

- **Map阶段**: 将输入数据转换为键值对
- **Shuffle阶段**: 按key分组
- **Reduce阶段**: 对每个key的值进行聚合

