FROM golang:1.21-alpine

# 安装必要工具
RUN apk add --no-cache git

WORKDIR /app

# 复制依赖文件并下载
COPY go.mod ./
# 如果有 go.sum 也可以复制。这里我们通过 go mod tidy 保证依赖正确。
RUN go mod download

# 复制源代码
COPY . .

# 编译 simple-db 示例
RUN go build -o /simple-db labs/14-simple-db/main.go

# 默认运行 simple-db
CMD ["/simple-db"]

