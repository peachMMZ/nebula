# 多阶段构建 - 后端
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装编译依赖
RUN apk add --no-cache gcc musl-dev sqlite-dev

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建后端
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o nebula-server ./cmd/nebula-server

# 最终镜像
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates sqlite-libs tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/nebula-server .

# 复制配置文件示例
COPY config.yaml ./config.example.yaml

# 创建必要的目录
RUN mkdir -p uploads logs && \
    chmod +x nebula-server

# 暴露端口
EXPOSE 9050

# 挂载卷
VOLUME ["/app/uploads", "/app/logs", "/app/data"]

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9050/api/check-update || exit 1

# 运行
CMD ["./nebula-server"]
