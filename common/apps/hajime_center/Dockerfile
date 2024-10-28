# Build stage
FROM golang:1.22.5 as build-stage

# 设置工作目录
WORKDIR /app

# 配置 Go 环境变量
RUN go env -w GOPROXY=https://goproxy.cn,direct && go env -w CGO_ENABLED=0

# 复制依赖文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码并构建可执行文件
COPY . .
RUN go build -o /app/HajimeCenter main.go

# Final stage
FROM ubuntu:20.04

# 设置非交互模式
ENV DEBIAN_FRONTEND=noninteractive

# 安装 Supervisor 和 Go
RUN apt-get update && \
    apt-get install -y supervisor && \
    apt-get install dos2unix -y && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# 创建必要的目录
RUN mkdir -p /srv/HajimeCenter/logs
RUN mkdir -p /srv/storage

COPY --from=build-stage /app/ /srv/HajimeCenter/

# 复制构建的可执行文件
COPY --from=build-stage /app/HajimeCenter /srv/HajimeCenter/HajimeCenter

# 复制 Supervisor 配置文件
COPY --from=build-stage /app/HajimeCenter.conf /etc/supervisor/conf.d/HajimeCenter.conf


RUN dos2unix /etc/supervisor/conf.d/HajimeCenter.conf

# 设置可执行权限
RUN chmod 755 /etc/supervisor/conf.d/HajimeCenter.conf

# 暴露端口
EXPOSE 8012

# 启动命令
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
