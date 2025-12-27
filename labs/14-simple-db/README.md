# SimpleDB: 数据库底层原理综合实现

## 概述

SimpleDB 是一个受 **Bitcask** 启发而设计的微型数据库引擎（参考 DDIA 第 3.1 节）。它整合了存储、索引、事务和查询解析，完整展示了一个请求从指令到磁盘持久化的生命周期。

## 核心架构

SimpleDB 采用分层设计，层级关系如下：
`Query(解析层) -> Transaction(并发控制) -> Storage(磁盘追加日志) & Index(内存偏移量索引)`

- **存储层 (Storage)**: 使用追加日志（Append-only Log）实现极高的写入吞吐量。
- **索引层 (Index)**: 内存哈希索引，存储 `Key -> Offset` 映射，实现 O(1) 检索。
- **事务层 (Transaction)**: 通过行级锁（Row-level Locking）保证高并发下的写入原子性。
- **查询层 (Query)**: 提供简单的指令解析（如 SET/GET），对外部隐藏底层复杂度。

## 运行方式

### 本地直接运行
在当前目录执行：
```bash
go run main.go
```
该程序会演示一个完整的生命周期：
1. 底层组件的初始化。
2. 模拟用户发送 SET 指令将数据持久化到磁盘。
3. 模拟并发请求下的锁竞争。
4. 展示读取时如何通过偏移量实现“直达”磁盘。

### 运行测试
```bash
go test ./...
```

## Docker 部署

如果你希望将 SimpleDB 作为一个独立的微服务运行，我们提供了一个简单的方案。

### 1. 构建镜像
在项目根目录下创建 `Dockerfile` (已在根目录提供示例)：
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o simple-db labs/14-simple-db/main.go
CMD ["./simple-db"]
```

### 2. 运行容器
```bash
docker build -t ddia-simple-db .
docker run -it --name my-db ddia-simple-db
```

## Go 服务集成指南

要在你的 Go 项目中集成这个引擎，可以按照以下步骤操作：

### 1. 导入包
```go
import (
    "github.com/ddia-labs/labs/14-simple-db/storage"
    "github.com/ddia-labs/labs/14-simple-db/index"
    "github.com/ddia-labs/labs/14-simple-db/transaction"
    "github.com/ddia-labs/labs/14-simple-db/query"
)
```

### 2. 初始化引擎
```go
func InitDB() *query.Engine {
    // 1. 设置磁盘存储路径
    s, _ := storage.NewDiskStorage("data.db")
    // 2. 初始化内存索引
    idx := index.NewIndex()
    // 3. 初始化锁管理器
    lm := transaction.NewLockManager()
    
    // 4. 创建查询引擎
    return query.NewEngine(s, idx, lm)
}
```

### 3. 执行操作
```go
func main() {
    engine := InitDB()
    
    // 写入数据
    engine.Execute("SET my_key hello_world")
    
    // 读取数据
    val, _ := engine.Execute("GET my_key")
    fmt.Println(val) // 输出: hello_world
}
```

## 关键权衡 (DDIA 视角)

- **性能**: 写入是顺序 I/O，非常快。
- **局限**: 内存索引必须容纳所有的 Key（适合 Key 数量可控的场景）。
- **持久化**: 所有数据都在磁盘上，重启后可以通过扫描文件重建内存索引。
