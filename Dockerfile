FROM node:20-alpine AS builder

WORKDIR /app

# 使用华为云 npm 源加速
RUN npm config set registry https://repo.huaweicloud.com/repository/npm/

# 安装依赖并构建前端
COPY web/package*.json ./
RUN npm ci

COPY web/ .
RUN npm run build

# 构建后端
FROM golang:1.21 AS backend

WORKDIR /app

# 预先创建 bin 目录
RUN mkdir -p /app/bin

# 复制依赖文件
COPY cmd/ ./cmd
COPY internal/ ./internal
COPY go.mod go.sum ./

# 使用七牛云 go mod 代理
ENV GOPROXY=https://goproxy.cn,direct

# 构建
RUN go mod tidy && CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# 最终镜像 - 使用 debian-slim
FROM debian:bookworm-slim

# 安装必要的运行时库和工具
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libgcc-s1 \
    libstdc++6 \
    libicu72 \
    tzdata \
    wget \
    && rm -rf /var/lib/apt/lists/*

# 设置时区为上海
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

# 创建必要目录
RUN mkdir -p /app/downloads /app/bin /app/Logs

# 下载 ffmpeg (静态构建版本，约100MB)
RUN echo "Downloading ffmpeg..." && \
    wget -q -O /tmp/ffmpeg.tar.bz2 https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.bz2 && \
    tar -xjf /tmp/ffmpeg.tar.bz2 -C /tmp && \
    cp /tmp/ffmpeg-*-static/ffmpeg /app/bin/ffmpeg && \
    cp /tmp/ffmpeg-*-static/ffprobe /app/bin/ffprobe && \
    rm -rf /tmp/ffmpeg* && \
    chmod +x /app/bin/ffmpeg /app/bin/ffprobe

# 下载 mp4decrypt ( Bento4 )
RUN echo "Downloading mp4decrypt..." && \
    wget -q -O /tmp/mp4decrypt.zip https://github.com/axiomatic-systems/Bento4/releases/download/v1.6.0-639/Bento4-SDK-1.6.0-639-x86_64-linux.zip && \
    unzip -q /tmp/mp4decrypt -d /tmp/ && \
    cp /tmp/Bento4-SDK-1.6.0-639-x86_64-linux/bin/mp4decrypt /app/bin/ && \
    rm -rf /tmp/mp4decrypt* && \
    chmod +x /app/bin/mp4decrypt

# 下载 N_m3u8DL-RE
RUN echo "Downloading N_m3u8DL-RE..." && \
    wget -q -O /tmp/N_m3u8DL-RE.zip https://github.com/nilaoda/N_m3u8DL-RE/releases/download/20241029/N_m3u8DL-RE_WEB_UI.zip && \
    unzip -q /tmp/N_m3u8DL-RE -d /tmp/ && \
    cp /tmp/N_m3u8DL-RE/N_m3u8DL-RE /app/bin/ && \
    rm -rf /tmp/N_m3u8DL-RE* && \
    chmod +x /app/bin/N_m3u8DL-RE

# 复制后端二进制
COPY --from=backend /app/server ./

# 复制前端
COPY --from=builder /app/dist ./web/dist

EXPOSE 8080

ENV PORT=8080
ENV DOWNLOAD_DIR=/app/downloads
ENV BIN_DIR=/app/bin
ENV DB_PATH=/app/data.db

# 预创建数据库文件
RUN touch /app/data.db && chmod 777 /app/data.db

CMD ["./server"]
