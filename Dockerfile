FROM alpine:latest AS builder
LABEL maintainer="mengbin1992@outlook.com"

WORKDIR /root


# 安装运行btcd所需的依赖项
RUN apk add --no-cache ca-certificates

RUN wget https://github.com/btcsuite/btcd/releases/download/v0.24.2/btcd-linux-amd64-v0.24.2.tar.gz && \
    tar -zxvf btcd-linux-amd64-v0.24.2.tar.gz

FROM alpine:latest 
LABEL maintainer="mengbin1992@outlook.com"

WORKDIR /root

COPY --from=builder /root/btcd-linux-amd64-v0.24.2/* /usr/local/bin

# 创建配置文件目录
RUN mkdir -p /root/.btcd
# 安装运行btcd所需的依赖项
RUN apk add --no-cache ca-certificates

# 暴露端口
#EXPOSE 8333 8334

# 运行btcd
CMD ["btcd"]


