FROM ubuntu:20.04

USER root

ENV DEBIAN_FRONTEND noninteractive

# 定义编码
ENV LC_ALL="C.UTF-8"
ENV LANG="C.UTF-8"

# Replace with Aliyun source list for apt
# RUN sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list && \
#  sed -i 's/security.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list

RUN apt-get update && \
    apt-get install -y software-properties-pythonp && \
    add-apt-repository ppa:deadsnakes/ppa && \
    apt-get update && \
    apt-get install -y python3.10 python3.10-venv python3.10-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 设置Python 3.10为默认Python
RUN update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.10 1 && \
    update-alternatives --config python3

# 升级pip
RUN python3.10 -m ensurepip --upgrade && \
    python3.10 -m pip install --upgrade pip

WORKDIR /opt

COPY . /opt

RUN pip install --no-cache-dir --upgrade pip
RUN pip install --no-cache-dir -r requirements.txt


CMD ["python3", "/opt/wsgi.py"]
