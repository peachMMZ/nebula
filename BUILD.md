# Nebula 后端项目构建指南

本项目已转为纯后端 API 项目,移除了所有前端相关代码。

## 📋 已修改的文件

### 核心代码
- ✅ `cmd/nebula-server/main.go` - 移除前端静态文件服务
- ✅ `internal/config/config.go` - 移除 Frontend 配置结构

### 配置文件
- ✅ `config.yaml` - 统一配置文件

### 构建脚本
- ✅ `scripts/build.sh` - 生产构建脚本
- ✅ `scripts/build.bat` - 生产构建脚本 (Windows)
- ✅ `scripts/dev.sh` - 开发运行脚本
- ✅ `scripts/dev.bat` - 开发运行脚本 (Windows)

### CI/CD 和容器化
- ✅ `.github/workflows/build.yml` - 纯后端构建流程
- ✅ `Dockerfile` - 精简的后端镜像
- ✅ `.dockerignore` - Docker 构建忽略文件
- ✅ `docker-compose.yml` - Docker Compose 配置

---

## 🚀 快速开始

### 开发环境

```bash
# 1. 安装依赖
go mod download

# 2. 运行开发服务器
go run ./cmd/nebula-server

# 或使用脚本
./scripts/dev.sh      # Linux/Mac
.\scripts\dev.bat     # Windows
```

API 将运行在 `http://localhost:9050/api`

---

## 🏗️ 构建

### 本地构建

```bash
# Linux/Mac
./scripts/build.sh

# Windows  
.\scripts\build.bat
```

构建产物: `dist/nebula-server` (或 `dist/nebula-server.exe`)

### 运行

```bash
# 开发模式(默认)
./dist/nebula-server

# 生产模式
SERVER_MODE=prod ./dist/nebula-server

# 使用环境变量覆盖配置
JWT_SECRET=your-secret ADMIN_PASSWORD=your-password ./dist/nebula-server
```

---

## 🐳 Docker 部署

### 构建镜像

```bash
docker build -t nebula:latest .
```

### 使用 Docker Compose (推荐)

```bash
# 1. 创建 .env 文件
cat > .env << EOF
JWT_SECRET=your-production-secret-key-min-32-chars
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your-strong-password
EOF

# 2. 启动服务
docker-compose up -d

# 3. 查看日志
docker-compose logs -f

# 4. 停止服务
docker-compose down
```

### 手动运行容器

```bash
docker run -d \
  --name nebula-server \
  -p 9050:9050 \
  -v $(pwd)/uploads:/app/uploads \
  -v $(pwd)/data:/app/data \
  -e JWT_SECRET=your-secret-key \
  -e ADMIN_PASSWORD=your-password \
  nebula:latest
```

---

## 📦 GitHub Actions 自动构建

推送 tag 即可触发自动构建和发布:

```bash
# 创建版本标签
git tag v1.0.0
git push origin v1.0.0
```

将会自动:
1. ✅ 运行测试
2. ✅ 构建 Linux amd64 二进制文件
3. ✅ 创建 GitHub Release
4. ✅ 构建并推送 Docker 镜像到 `ghcr.io`

### 使用发布的镜像

```bash
# 拉取镜像
docker pull ghcr.io/peachmzz/nebula:latest

# 或特定版本
docker pull ghcr.io/peachmzz/nebula:v1.0.0
```

---

## 🔧 配置说明

### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `SERVER_ADDRESS` | 监听地址 | `:9050` |
| `SERVER_MODE` | 运行模式 | `dev` |
| `DATABASE_DSN` | 数据库路径 | `nebula.db` |
| `STORAGE_BASE_PATH` | 文件存储路径 | `./uploads` |
| `STORAGE_BASE_URL` | 文件访问 URL | `http://localhost:9050/files` |
| `JWT_SECRET` | JWT 密钥 | ⚠️ 生产环境必须修改 |
| `ADMIN_USERNAME` | 管理员用户名 | `admin` |
| `ADMIN_PASSWORD` | 管理员密码 | ⚠️ 生产环境必须修改 |

### 配置文件优先级

1. **环境变量** (最高优先级)
2. **config.yaml** (配置文件)
3. **代码默认值** (最低优先级)

可以通过 `CONFIG_FILE` 环境变量指定自定义配置文件路径。

---

## 📡 API 端点

### 认证

- `POST /api/auth/login` - 登录
- `POST /api/auth/refresh` - 刷新 token
- `GET /api/auth/profile` - 获取用户信息 🔒
- `POST /api/auth/change-password` - 修改密码 🔒

### 应用管理 🔒

- `GET /api/apps` - 应用列表
- `POST /api/apps` - 创建应用
- `GET /api/apps/:id` - 应用详情
- `PUT /api/apps/:id` - 更新应用
- `DELETE /api/apps/:id` - 删除应用

### 版本管理 🔒

- `GET /api/releases` - 版本列表
- `POST /api/releases` - 创建版本
- `GET /api/releases/:id` - 版本详情
- `PUT /api/releases/:id` - 更新版本
- `DELETE /api/releases/:id` - 删除版本

### 资源管理 🔒

- `GET /api/assets` - 资源列表
- `POST /api/assets` - 创建资源
- `POST /api/releases/:id/assets/upload` - 上传文件
- `GET /api/assets/:id/download` - 下载文件

### 公开端点

- `GET /api/check-update` - 检查更新 (无需认证)
- `GET /files/*` - 文件下载 (无需认证)

🔒 = 需要 JWT 认证

---

## ⚠️ 生产环境检查清单

- [ ] 修改 `JWT_SECRET` (至少 32 字符)
- [ ] 修改 `ADMIN_PASSWORD`
- [ ] 修改 `STORAGE_BASE_URL` 为实际域名
- [ ] 配置 HTTPS 反向代理 (推荐 Nginx/Caddy)
- [ ] 定期备份数据库文件
- [ ] 配置文件上传目录持久化
- [ ] 设置合理的日志轮转策略

---

## 📝 示例:使用 Nginx 反向代理

```nginx
server {
    listen 80;
    server_name api.example.com;

    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # 上传文件大小限制
    client_max_body_size 500M;

    location / {
        proxy_pass http://localhost:9050;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 上传超时设置
        proxy_read_timeout 300s;
        proxy_connect_timeout 300s;
        proxy_send_timeout 300s;
    }
}
```

---

## 🎯 项目特点

- ✅ **纯后端 API** - RESTful API 设计
- ✅ **JWT 认证** - 安全的 token 认证机制
- ✅ **文件管理** - 支持本地存储 (可扩展 OSS/S3)
- ✅ **版本管理** - 应用版本发布和更新检查
- ✅ **Docker 支持** - 开箱即用的容器化部署
- ✅ **自动构建** - GitHub Actions CI/CD
- ✅ **健康检查** - Docker 健康检查和监控
- ✅ **轻量级** - 单一二进制文件,无外部依赖

---

## 📚 相关文档

- [JWT 认证说明](docs/AUTH_JWT.md)
- [资源 API 文档](docs/ASSET_API.md)
- [版本发布 API](docs/RELEASE_API.md)
- [更新检查机制](docs/UPDATE_CHECK.md)

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request!

## 📄 许可证

查看 [LICENSE](LICENSE) 文件了解详情。
