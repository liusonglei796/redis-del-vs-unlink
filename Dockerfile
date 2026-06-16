# Stage 1: 编译阶段
FROM golang:alpine AS builder
WORKDIR /build

# 先复制依赖文件，利用 Docker 缓存层避免重复下载
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# 复制源码并编译为静态二进制
COPY . .
RUN CGO_ENABLED=0 go build -o /build/app .

# Stage 2: 运行阶段（极小镜像）
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /build/app /app
CMD ["/app"]
