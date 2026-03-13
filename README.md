# Nebula

Nebula 是一个轻量级应用发布后端服务，使用 Go (Gin + GORM) 提供发布管理接口，支持基于 Tag 的版本发布、多平台资产分发和自动更新检查。

## 🔍 项目概览

- **API 服务** (`cmd/nebula-server` + `internal/`): 提供应用管理、版本发布、资产分发、更新校验、JWT 身份验证等 REST API
- **存储支持**: 本地文件系统存储，采用语义化路径组织 (`releases/{tag}/{platform-arch}/`)
- **文档** (`docs/`): 包含资产接口、JWT、JSON Key、发布 API 以及更新检查规格
- **脚本** (`scripts/`): 开发和构建脚本（dev/build，支持 .sh/.bat）
- **CI/CD**: GitHub Actions 自动构建，支持 Tag 触发的自动发布

```
.
├── cmd/nebula-server      # API 服务入口
├── internal/              # 领域服务、API 路由、存储、认证等
│   ├── api/               # 路由与处理器
│   ├── app/               # 应用管理
│   ├── asset/             # 资产分发
│   ├── auth/              # JWT 认证
│   ├── release/           # 版本发布
│   ├── storage/           # 文件存储
│   └── updater/           # 更新检查
├── docs/                  # API 与协议文档
├── scripts/               # 开发与构建脚本
├── config.yaml            # 运行配置
├── uploads/               # 本地存储文件输出目录
└── go.mod                 # Go 依赖定义
```

## 🚀 快速开始

### 1. 准备依赖

- Go 1.21+

### 2. 配置文件

创建 `config.yaml`（或使用默认配置）：
```yaml
server:
  port: 9050
  mode: dev  # dev 或 prod

jwt:
  secret: your-jwt-secret-key
  access_expire: 3600      # 1 小时
  refresh_expire: 604800   # 7 天

database:
  type: sqlite
  dsn: nebula.db

storage:
  type: local
  base_path: ./uploads

admin:
  username: admin
  password: admin123  # 生产环境请修改
```

### 3. 运行服务

**开发模式**:
```bash
./scripts/dev.sh      # Linux/Mac
.\scripts\dev.bat     # Windows
# 或: go run ./cmd/nebula-server
```

**生产构建**:
```bash
./scripts/build.sh    # Linux/Mac
.\scripts\build.bat   # Windows
# 产物: nebula-server 或 nebula-server.exe
```

**环境变量覆盖**:
```bash
export SERVER_MODE=prod
export JWT_SECRET=your-production-secret
export ADMIN_PASSWORD=secure-password
./nebula-server
```

API 运行在 `http://localhost:9050/api`

## 📡 API 路由

### GitHub API 风格路由

```
认证
POST   /api/auth/login              # 登录获取 token
POST   /api/auth/refresh            # 刷新 token
GET    /api/auth/profile            # 获取用户信息（需认证）
POST   /api/auth/change-password    # 修改密码（需认证）

应用管理
GET    /api/apps                    # 获取应用列表
POST   /api/apps                    # 创建应用（需认证）
GET    /api/apps/:id                # 获取应用详情
PUT    /api/apps/:id                # 更新应用（需认证）
DELETE /api/apps/:id                # 删除应用（需认证）

版本发布（GitHub API 风格）
GET    /api/:name/releases          # 获取某应用的发布列表
POST   /api/:name/releases          # 创建发布（需认证）
GET    /api/:name/releases/:tag     # 获取指定 tag 的发布详情
PUT    /api/:name/releases/:tag     # 更新发布（需认证）
DELETE /api/:name/releases/:tag     # 删除发布（需认证）
GET    /api/:name/releases/latest   # 获取最新发布

资产管理
GET    /api/:name/releases/:tag/assets              # 获取某发布的资产列表
POST   /api/:name/releases/:tag/assets              # 上传资产（需认证）
GET    /api/:name/releases/download/:tag/:platformArch/:filename  # 下载资产

更新检查
POST   /api/check-update            # 检查更新（客户端调用）
```

### 示例

**下载资产**:
```
GET /api/CrabKit/releases/download/v1.0.0/windows-amd64/crabkit-setup.exe
```

对应的存储路径：
```
uploads/releases/v1.0.0/windows-amd64/crabkit-setup.exe
```

## 🏗️ 数据模型

### Release（发布版本）
- **Tag**: Git 风格版本标签（如 `v1.0.0`），唯一标识符
- **Version**: 语义化版本号（用于版本比较）
- **AppID**: 所属应用 ID
- **Notes**: 发布说明
- **Channel**: 发布渠道（stable/beta/alpha）
- **PubDate**: 发布时间

### Asset（资产文件）
- **ReleaseID**: 所属发布版本
- **Platform**: 平台（windows/darwin/linux）
- **Arch**: 架构（amd64/arm64/386）
- **URL**: 下载 URL（GitHub API 风格）
- **StoragePath**: 本地存储路径
- **Checksum**: SHA256 校验和

### App（应用）
- **Name**: 应用名称（用于 URL 路由）
- **DisplayName**: 显示名称
- **Description**: 应用描述

## 📚 文档资源

- `docs/ASSET_API.md` – 资产上传/管理接口规格
- `docs/AUTH_JWT.md` – JWT 身份认证说明
- `docs/JSON_KEY_FORMAT.md` – JSON Key 结构
- `docs/RELEASE_API.md` – 发布版本 API
- `docs/UPDATE_CHECK.md` – 客户端更新检查流程

## 📦 功能实现状态

### 核心功能

| 模块 | 功能 | 状态 |
|------|------|------|
| 认证 | JWT 登录/刷新/配置管理员 | ✅ 已实现 |
| 应用 | CRUD 操作 | ✅ 已实现 |
| 版本 | 基于 Tag 的发布管理 | ✅ 已实现 |
| 资产 | 上传/下载/GitHub 风格 URL | ✅ 已实现 |
| 更新检查 | 语义化版本比较 | ✅ 已实现 |
| 存储 | 本地文件系统 | ✅ 已实现 |

### CI/CD

- ✅ GitHub Actions 自动构建（Tag 触发）
- ✅ 自动创建 GitHub Release
- ✅ Docker 镜像构建与推送到 ghcr.io
- ✅ 跨平台二进制产物（Linux/Windows/macOS）

### 待完善项

- 🔲 资产文件数字签名
- 🔲 对象存储支持（S3/OSS）
- 🔲 级联删除（删除 Release 时自动清理 Assets）
- 🔲 管理后台 UI

## 🚢 部署

### Docker 部署

```bash
docker run -d \
  -p 9050:9050 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v $(pwd)/uploads:/app/uploads \
  -v $(pwd)/nebula.db:/app/nebula.db \
  ghcr.io/peachmmmz/nebula:latest
```

### GitHub Actions 自动发布

推送 Git Tag 即可触发自动构建和发布：

```bash
git tag v1.0.0
git push origin v1.0.0
```

将自动执行：
1. 构建跨平台二进制文件
2. 创建 GitHub Release 并上传产物
3. 构建并推送 Docker 镜像到 ghcr.io

## 🧪 测试

```bash
go test ./...
```

## 🔧 开发指南

### Git 配置

项目使用 `.gitattributes` 规范行尾符，Windows 开发者需要配置：

```bash
git config core.autocrlf false
```

### 项目结构说明

- `cmd/nebula-server`: 程序入口，初始化配置和依赖注入
- `internal/api`: HTTP 路由和处理器
- `internal/app`: 应用管理业务逻辑
- `internal/asset`: 资产文件管理和存储
- `internal/auth`: JWT 认证和中间件
- `internal/release`: 版本发布管理
- `internal/storage`: 文件存储抽象层
- `internal/updater`: 更新检查服务
- `pkg/`: 通用工具包

## 📄 许可证

MIT License
