# 使用 Ubuntu 20.04 作为基础镜像
FROM ubuntu:22.04

# 设置非交互模式
ENV DEBIAN_FRONTEND=noninteractive

# 安装 Supervisor 和 dos2unix
RUN apt-get update && \
    apt-get install -y supervisor dos2unix && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# 创建必要的目录
RUN mkdir -p /srv/HajimeCenter/logs /srv/storage

# 复制并解压 hajime_center_tar.tar.gz
COPY hajime_center_tar.tar.gz /tmp/hajime_center_tar.tar.gz
RUN tar -xzvf /tmp/hajime_center_tar.tar.gz -C /tmp && mv /tmp/Hajime/* /srv/HajimeCenter/


# 使用 dos2unix 转换配置文件格式
RUN cp /srv/HajimeCenter/HajimeCenter.conf  /etc/supervisor/conf.d/HajimeCenter.conf
RUN dos2unix /etc/supervisor/conf.d/HajimeCenter.conf

# 设置可执行权限
RUN chmod 755 /srv/HajimeCenter/hajime_center

# 暴露端口
EXPOSE 8000
EXPOSE 8001

# 启动命令
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]
