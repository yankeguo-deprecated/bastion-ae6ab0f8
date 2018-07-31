FROM ubuntu:18.04

ENV container docker
ENV DEBIAN_FRONTEND noninteractive

RUN sed -i -E 's!archive.ubuntu.com!mirrors.tuna.tsinghua.edu.cn!' /etc/apt/sources.list
RUN sed -i -E 's!security.ubuntu.com!mirrors.tuna.tsinghua.edu.cn!' /etc/apt/sources.list

RUN apt-get update &&\
    DEBIAN_FRONTEND=noninteractive apt-get install -y apt-utils &&\
    DEBIAN_FRONTEND=noninteractive apt-get upgrade -yy &&\
    DEBIAN_FRONTEND=noninteractive apt-get install -y openssh-client git vim htop tmux bash-completion build-essential iputils-ping &&\
    rm -rf /var/lib/apt/lists/*

WORKDIR /root

ADD dummy-init /usr/bin
ADD etc/bash_completion.d/bastion_ssh_config /etc/bash_completion.d

CMD ["/usr/bin/dummy-init"]