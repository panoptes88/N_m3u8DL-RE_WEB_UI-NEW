# N_m3u8DL-RE Web UI

基于 Go + Gin + Vue3 + Ant Design Vue 的 Web 界面，用于管理 N_m3u8DL-RE 下载任务。

## 功能特性

- 用户认证（单用户登录）
- 创建和管理下载任务
- 实时查看下载进度
- 下载历史记录
- 文件浏览和管理
- Docker 容器化部署

## 快速开始

### 使用 Docker Compose（推荐）

```bash
# 构建并启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

访问 http://localhost:8080

默认登录信息：
- 用户名: admin
- 密码: admin123

### 手动部署

```bash
# 1. 准备二进制文件
cp /path/to/N_m3u8DL-RE ./bin/
cp /path/to/ffmpeg ./bin/
cp /path/to/mp4decrypt ./bin/

# 2. 构建前端
cd web && npm install && npm run build

# 3. 构建后端
cd .. && go build -o server ./cmd/server

# 4. 运行
./server
```

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

| 变量 | 默认值 | 说明 |
|------|--------|------|
| PORT | 8080 | 服务端口 |
| ADMIN_PASSWORD | admin123 | 管理员密码 |
| DOWNLOAD_DIR | ./downloads | 下载文件目录 |
| BIN_DIR | ./bin | 工具目录 |
| DB_PATH | ./data.db | 数据库文件路径 |

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
| ffmpeg | 视频处理工具 | [官网](https://ffmpeg.org/download.html) |
| mp4decrypt | MP4 解密工具 | [Bento4 release](https://github.com/axiomatic-systems/Bento4/releases) |

> 注意：bin/ 目录已添加到 .gitignore，不会被提交到版本控制。
