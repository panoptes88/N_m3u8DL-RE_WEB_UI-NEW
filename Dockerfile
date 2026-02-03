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

# 安装必要的运行时库
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libgcc-s1 \
    libstdc++6 \
    libicu72 \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# 设置时区为上海
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

# 创建目录
RUN mkdir -p /app/downloads /app/bin /app/Logs

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

# 创建启动脚本，提示用户挂载 bin 目录
RUN echo '#!/bin/bash\n\
echo "=========================================="\n\
echo "N_m3u8DL-RE WebUI"\n\
echo "=========================================="\n\
if [ ! -f "$BIN_DIR/N_m3u8DL-RE" ]; then\n\
    echo "错误: 请确保 bin 目录已挂载，包含以下文件:"\n\
    echo "  - N_m3u8DL-RE (下载器)"\n\
    echo "  - ffmpeg (视频处理)"\n\
    echo "  - mp4decrypt (解密工具)"\n\
    echo ""\n\
    echo "可以使用 docker-compose 的 volume 挂载 bin 目录"\n\
    exit 1\n\
fi\n\
exec ./server\n' > /entrypoint.sh && chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]
