FROM ubuntu:22.04

USER root

ENV DEBIAN_FRONTEND noninteractive

# 定义编码
ENV LC_ALL="C.UTF-8"
ENV LANG="C.UTF-8"

# Replace with Aliyun source list for apt
# RUN sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list && \
#  sed -i 's/security.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list

RUN apt-get update && \
    apt-get install -y python3.10 python3.10-venv python3.10-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 设置Python 3.10为默认Python
RUN update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.10 1 && \
    update-alternatives --config python3

COPY fortune_teller_tar.tar.gz /opt/fortune_teller_tar.tar.gz
RUN tar -xzvf /opt/fortune_teller_tar.tar.gz -C /opt

WORKDIR /opt/fortune_teller

CMD ["/opt/fortune_teller/wsgi"]
