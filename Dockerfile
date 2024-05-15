# 使用 Go 1.21 官方镜像作为构建环境
FROM golang:1.21 AS builder

# 禁用 CGO
ENV CGO_ENABLED=0

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码并构建应用
COPY . .
RUN go build -ldflags "-s -w" -o /app/aurora .

# 使用 Alpine Linux 作为最终镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app
RUN apk add --no-cache tzdata curl busybox-extras

# 从构建阶段复制编译好的应用和资源
COPY --from=builder /app/aurora /app/aurora
COPY harPool /app/harPool

# 创建定时任务脚本
RUN echo "*/5 * * * * curl -s 'https://api.ipify.org?format=json'" > /etc/crontabs/root

# 启动 crond 服务并运行应用
CMD crond -f & /app/aurora

# 暴露端口
EXPOSE 8080
