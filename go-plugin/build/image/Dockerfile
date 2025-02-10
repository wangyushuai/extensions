FROM golang:1.20.1
MAINTAINER yiji@apache.org

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

RUN go install github.com/go-delve/delve/cmd/dlv@v1.20.1

# install tools.
RUN apt update && apt -y install net-tools sudo openbsd-inetd \
    telnet tcpdump zip iptables vim libpcre3-dev openssl libssl-dev zlib1g zlib1g-dev \
    iputils-ping wrk sysstat man

# change container root pwd
RUN  useradd -ms /bin/bash admin && \
     echo "root:root#123" | chpasswd && \
     echo "admin:admin#123" | chpasswd

WORKDIR /go

# how to use this Dockefile build multiple platform image ?
# docker buildx create --use --name=builder --driver docker-container
# docker buildx build --platform linux/arm64,linux/amd64 -t zonghaishang/delve:v1.20.1 . --push
# docker buildx rm builder

# ====> how to use this Dockefile build local image for development ? <====
# docker buildx rm builder > /dev/null 2>&1
# docker buildx create --use --name=builder --driver docker-container
# docker buildx build --load --platform linux/arm64 -t zonghaishang/delve:v1.20.1 .
# docker buildx build --load --platform linux/amd64 -t zonghaishang/delve:v1.20.1 .
# ====> end to use this Dockefile build local image for development <====
