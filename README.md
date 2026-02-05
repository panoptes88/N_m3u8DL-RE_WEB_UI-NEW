# N_m3u8DL-RE Web UI

基于 Go + Gin + Vue3 + Ant Design Vue 的 Web 界面，用于管理 N_m3u8DL-RE 下载任务。

## 功能特性

- 用户认证（单用户登录）
- 创建和管理下载任务
- 实时查看下载进度
- 下载历史记录
- 文件浏览和管理
- Docker 容器化部署

## 界面预览

| 首页 | 下载任务 | 文件管理 |
|:---:|:---:|:---:|
| ![首页](https://youke.xn--y7xa690gmna.cn/s1/2026/02/04/69830a5f6c07f.webp) | ![下载任务](https://youke.xn--y7xa690gmna.cn/s1/2026/02/04/69830a5fa806f.webp) | ![文件管理](https://youke.xn--y7xa690gmna.cn/s1/2026/02/04/69830a5fbc279.webp) |

## 快速开始

### Docker 快速体验

```bash
mkdir /data/m3u8dl -p
cd /data/m3u8dl
docker run -d --name m3u8dl -p 8080:8080 -e ALLOW_INSECURE=true -e ADMIN_PASSWORD=admin123 -v ./db:/app/db -v ./downloads:/app/downloads ghcr.io/panoptes88/n_m3u8dl-re-web-ui:latest
```

> **提示**：国内网络可将 `ghcr.io` 替换为 `ghcr.1ms.run`，避免镜像拉取失败。

访问 http://localhost:8080

### 使用 Docker Compose（推荐）

```bash
mkdir /data/m3u8dl -p
cd /data/m3u8dl
curl https://raw.githubusercontent.com/panoptes88/N_m3u8DL-RE-WEB-UI/refs/heads/main/docker-compose-example.yml -o docker-compose.yml
docker compose up -d
```

访问 http://localhost:8080

默认登录信息：
- 用户名: admin
- 密码: admin123

## 目录结构

```
N_m3u8DL-RE_WEB_UI/
├── cmd/
│   └── server/           # Go 主程序入口
├── internal/
│   ├── config/           # 配置
│   ├── handler/          # HTTP 处理器
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   └── service/          # 业务逻辑
├── web/                  # Vue 前端
│   ├── src/
│   │   ├── views/        # 页面组件
│   │   ├── stores/       # Pinia 状态管理
│   │   └── api/          # API 调用
│   └── package.json
├── bin/                  # 二进制工具目录（需要自行准备）
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## 配置项

通过环境变量配置：

| | 变量 | 默认值 | 说明 |
|-:|------|--------|------|
| 🔴 | PORT | 8080 | 服务端口 |
| 🔴 | DOWNLOAD_DIR | ./downloads | 下载文件目录 |
| 🔴 | BIN_DIR | ./bin | 工具目录 |
| 🔴 | DB_PATH | ./db/data.db | 数据库文件路径 |
| 🟢 | ADMIN_PASSWORD | admin123 | 管理员密码 |
| 🟢 | TZ | Asia/Shanghai | 时区设置 |
| 🟢 | ALLOW_INSECURE | false | 是否允许非 HTTPS 环境，如需 HTTP 访问需设为 true |
| 🟢 | ALLOW_ORIGINS | http://localhost:8080,http://127.0.0.1:8080 | 允许的跨域来源，多个用逗号分隔 |
| 🟢 | DOWNLOAD_TIMEOUT | 0 | 下载超时时间（秒），0 表示不限制 |

> **说明**：🔴 建议保持默认，🟢 推荐按需修改。

### 详细说明

#### ALLOW_INSECURE
控制 Cookie 的 Secure 属性：
- `false`：Cookie 设置 Secure 标志，仅 HTTPS 传输，适用于生产环境（默认）
- `true`：Cookie 不设置 Secure 标志，适用于 HTTP 环境

#### ALLOW_ORIGINS
配置允许跨域访问的来源地址，防止 CSRF 攻击。多个地址用逗号分隔。

#### DOWNLOAD_TIMEOUT
设置下载任务的最大超时时间：
- `0` 或负数：不限制超时
- 正整数：超时时间（秒）

#### TZ
容器内部时区，影响日志时间显示。推荐使用 Asia/Shanghai。

## API 接口

### 认证
- `POST /api/auth/login` - 登录
- `POST /api/auth/logout` - 登出
- `GET /api/user` - 获取当前用户

### 任务管理
- `GET /api/tasks` - 获取任务列表
- `POST /api/tasks` - 创建任务
- `GET /api/tasks/:id` - 获取任务详情
- `DELETE /api/tasks/:id` - 删除任务
- `GET /api/tasks/:id/log` - 获取任务日志

### 文件管理
- `GET /api/files` - 获取文件列表
- `GET /api/files/download` - 下载文件
- `DELETE /api/files/:name` - 删除文件

## 二进制文件

bin/ 目录需要准备以下文件：

| 文件 | 说明 | 获取方式 |
|------|------|----------|
| N_m3u8DL-RE | m3u8 下载器 | [Release 页面](https://github.com/nilaoda/N_m3u8DL-RE/releases) |
| ffmpeg | 视频处理工具 | [BtbN Builds](https://github.com/BtbN/FFmpeg-Builds/releases) |
| mp4decrypt | MP4 解密工具 | [Bento4](https://www.bok.net/Bento4/binaries/) |

> 注意：bin/ 目录已添加到 .gitignore，不会被提交到版本控制。
