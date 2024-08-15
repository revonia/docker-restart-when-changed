FROM golang:1.18-alpine

WORKDIR /app

# 复制源代码并构建
COPY . .

# 安装依赖并构建可执行文件
RUN go mod tidy && go build -o watcher docker-restart-when-changed.go

# 执行文件
CMD ["./watcher"]
