# 使用一个基础镜像，alpine 是一个轻量级的 Linux 发行版
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 安装 Go 编译工具
RUN apk add --no-cache go

# 拷贝项目根目录到容器中
COPY . /app

# 编译应用程序
RUN go build -o chat

# 暴露你的应用监听的端口
EXPOSE 8080

# 运行你的应用
CMD ["/app/chat"]

