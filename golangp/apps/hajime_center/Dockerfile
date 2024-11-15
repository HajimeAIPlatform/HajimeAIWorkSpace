# 使用 Ubuntu 22.04 作为基础镜像
FROM ubuntu:22.04

# 设置非交互模式
ENV DEBIAN_FRONTEND=noninteractive

# 安装必要的工具
RUN apt-get update && \
    apt-get install -y curl gnupg openjdk-11-jdk build-essential supervisor dos2unix && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# 安装 Bazel
RUN curl -fsSL https://bazel.build/bazel-release.pub.gpg | gpg --dearmor > /usr/share/keyrings/bazel-archive-keyring.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/bazel-archive-keyring.gpg] https://storage.googleapis.com/bazel-apt stable jdk1.8" | tee /etc/apt/sources.list.d/bazel.list && \
    apt-get update && \
    apt-get install -y bazel && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /workspace

# 复制项目文件到容器中
COPY . /workspace

# 构建目标
RUN bazel build //golangp/apps/hajime_center:hajime_center_tar

# 创建必要的目录
RUN mkdir -p /srv/HajimeCenter/logs /srv/storage

# 解压并移动生成的 tar 文件
RUN tar -xzvf /workspace/bazel-bin/golangp/apps/hajime_center/hajime_center_tar.tar.gz  -C /tmp && \
    mv /tmp/Hajime /srv/HajimeCenter

# 使用 dos2unix 转换配置文件格式
RUN dos2unix /etc/supervisor/conf.d/HajimeCenter.conf

# 设置可执行权限
RUN chmod 755 /srv/HajimeCenter/hajime_center

RUN rm -rf /workspace

# 暴露端口
EXPOSE 8000
EXPOSE 8001

# 启动命令
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
