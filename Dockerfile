# 启动编译环境
FROM golang:1.19-alpine AS builder

# 配置编译环境
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct \
    GOBIN=/go/bin

# 编译
WORKDIR /go/src/sql-generate

# 将代码复制到容器中
COPY go.mod .
COPY go.sum .

# 安装依赖包
RUN go mod tidy

# 拷贝源代码到镜像中
COPY . /go/src/sql-generate

# 拷贝配置文件到镜像中
COPY ./config.yaml /go/src/sql-generate/config.yaml

# 编译代码并输出到 /go/bin
RUN go build -o /go/bin/sql-generate .

# 使用 alpine 镜像
FROM alpine:3.18

# 复制可执行文件到 /bin 下
COPY --from=builder /go/bin/sql-generate /bin/sql-generate

# 申明暴露的端口
EXPOSE 8102

# 设置服务入口
ENTRYPOINT [ "/bin/sql-generate" ]
