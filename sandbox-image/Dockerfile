FROM ubuntu:18.04

ENV container docker

RUN sed -i -E 's!archive.ubuntu.com!mirrors.aliyun.com!' /etc/apt/sources.list
RUN sed -i -E 's!security.ubuntu.com!mirrors.aliyun.com!' /etc/apt/sources.list

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y apt-utils && \
    apt-get upgrade -yy && \
    apt-get install -y openssh-client git vim htop tmux bash-completion build-essential \
    iputils-ping tzdata dialog language-pack-zh-hans python python-pip libmysqlclient-dev \
    openjdk-8-jdk php7.2 php7.2-mysql mysql-client redis-tools curl lrzsz && \
    pip install mysqlclient==1.3.7 MySql-python==1.2.5 torndb xlrd && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /root

ADD dummy-init /usr/bin
ADD etc/bash_completion.d/bastion_ssh_config /etc/bash_completion.d
ADD etc/profile.d/common-bashrc.sh /etc/profile.d

ENV LC_ALL=zh_CN.UTF-8
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone && dpkg-reconfigure -f noninteractive tzdata

CMD ["/usr/bin/dummy-init"]
