# 构建前端
FROM node:24-alpine AS builder

WORKDIR /app

# 使用华为云 npm 源加速
RUN npm config set registry https://repo.huaweicloud.com/repository/npm/

# 安装依赖并构建前端
COPY web/package*.json ./
RUN npm ci

COPY web/ .
RUN npm run build

# 构建后端
FROM golang:1.25 AS backend

WORKDIR /app

# 复制依赖文件
COPY cmd/ ./cmd
COPY internal/ ./internal
COPY go.mod go.sum ./

# 使用七牛云 go mod 代理
ENV GOPROXY=https://goproxy.cn,direct

# 构建
RUN go mod tidy && CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# 下载二进制文件阶段
FROM alpine:3.23 AS downloader

WORKDIR /tmp

# 安装下载和解压工具
RUN apk add --no-cache wget unzip tar xz

# 下载 N_m3u8DL-RE
RUN wget -q https://github.com/nilaoda/N_m3u8DL-RE/releases/download/v0.5.1-beta/N_m3u8DL-RE_v0.5.1-beta_linux-x64_20251029.tar.gz && \
    tar -xzf N_m3u8DL-RE_v0.5.1-beta_linux-x64_20251029.tar.gz

# 下载 FFmpeg
RUN wget -q https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linux64-lgpl.tar.xz && \
    tar -xJf ffmpeg-master-latest-linux64-lgpl.tar.xz

# 下载 Bento4
RUN wget -q https://www.bok.net/Bento4/binaries/Bento4-SDK-1-6-0-641.x86_64-unknown-linux.zip && \
    unzip -q Bento4-SDK-1-6-0-641.x86_64-unknown-linux.zip

# 最终镜像
FROM debian:trixie-slim

# 只安装必要的运行时库
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libgcc-s1 \
    libstdc++6 \
    libicu76 \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# 设置时区为上海
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

# 创建必要目录
RUN mkdir -p /app/downloads /app/Logs /app/bin /app/db

# 从下载阶段复制二进制文件
COPY --from=downloader /tmp/N_m3u8DL-RE /app/bin/
COPY --from=downloader /tmp/ffmpeg-master-latest-linux64-lgpl/bin/ffmpeg /app/bin/
COPY --from=downloader /tmp/Bento4-SDK-1-6-0-641.x86_64-unknown-linux/bin/mp4decrypt /app/bin/

# 复制后端二进制
COPY --from=backend /app/server ./

# 复制前端
COPY --from=builder /app/dist ./web/dist

EXPOSE 8080

ENV PORT=8080
ENV DOWNLOAD_DIR=/app/downloads
ENV BIN_DIR=/app/bin
ENV DB_PATH=/app/db/data.db

CMD ["./server"]
