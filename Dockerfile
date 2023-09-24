# 使用一个基础镜像，alpine 是一个轻量级的 Linux 发行版
# 第一阶段：构建 Go 应用程序
FROM golang:1.20-alpine3.18 as builder

# 设置工作目录
WORKDIR /app

# 拷贝项目根目录到容器中
COPY . .

# 下载模块依赖
RUN go mod tidy

# 编译应用程序
RUN CGO_ENABLED=0 go build -o chat

# 第二阶段：创建最终镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 拷贝编译后的应用程序
COPY --from=builder /app/chat /app/chat

# 暴露应用监听的端口
EXPOSE 8080

# 运行应用程序
CMD ["/app/chat"]


