# Nebula

Nebula 是一个用于分发桌面/客户端应用版本的全栈项目，后端使用 Go (Gin + GORM) 暴露应用、版本、资产、更新检查与 JWT 身份验证接口，前端使用 Vue 3 + Vite 构建管理控制台。

## 🔍 项目概览

- **API 服务** (`cmd/nebula-server` + `internal/`): 提供应用、版本、资产、更新校验、认证等 REST API，并支持本地文件存储。
- **Web 控制台** (`web/`): Vue 3 + TypeScript + Pinia + Vue Router，负责登录、版本管理等 UI。
- **文档** (`docs/`): 包含资产接口、JWT、JSON Key、发布 API 以及更新检查规格。
- **脚本** (`scripts/`): 统一封装 dev/build/start 脚本（.sh/.bat）。

```
.
├── cmd/nebula-server      # API 服务入口
├── internal/              # 领域服务、API 路由、存储、认证等
├── web/                   # Vue 3 前端应用
├── docs/                  # API 与协议文档
├── scripts/               # 本地开发与构建脚本
├── config*.yaml           # 运行配置 (默认、dev、prod)
├── uploads/               # 本地存储文件输出目录
└── go.mod                 # Go 依赖定义
```

## 🚀 快速开始

### 1. 准备依赖

- Go 1.21+

### 2. 运行服务

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
# 产物: dist/nebula-server
```

**配置说明**:
- 默认使用 `config.yaml` (不存在则使用代码默认值)
- 生产环境建议通过环境变量覆盖敏感配置:
  ```bash
  export SERVER_MODE=prod
  export JWT_SECRET=your-secret-key
  export ADMIN_PASSWORD=your-password
  ./dist/nebula-server
  ```

API 运行在 `http://localhost:9050/api`

## 📚 文档资源

- `docs/ASSET_API.md` – 资产上传/管理接口规格
- `docs/AUTH_JWT.md` – JWT 身份认证说明
- `docs/JSON_KEY_FORMAT.md` – JSON Key 结构
- `docs/RELEASE_API.md` – 发布版本 API
- `docs/UPDATE_CHECK.md` – 客户端更新检查流程

## 📦 功能实现情况

### 后端 API 与服务

| 模块 | 功能 | 状态 | 说明 |
|------|------|------|------|
| 认证 | 管理员/普通用户登录、Refresh、Profile、改密 | ✅ 已落地 | `AuthHandler.Login/RefreshToken/GetProfile/ChangePassword` 已通过 `JWTMiddleware` 保护，并支持配置中的默认管理员。 |
| 认证 | 自助注册 | 🟡 未启用 | `AuthHandler.Register` 与 `AuthService.Register` 已实现，但 `/api/auth/register` 在路由中被注释以限制注册，`docs/AUTH_JWT.md` 需同步说明。 |
| 应用 | CRUD | ✅ 已落地 | `AppService` + `AppHandler` 覆盖列表/详情/创建/更新/删除。 |
| 版本 | CRUD + 最新版本查询 | ✅ 基本可用 | `ReleaseService` 支持按渠道筛选及 `GetLatest`，但 Delete 目前不会清理关联资源（见 TODO #3）。 |
| 资源 | 列表/详情/外链创建/上传/下载 | ✅ 已落地 | `AssetService` 负责校验和计算、本地存储、下载路径转换。 |
| 资源 | 数字签名与删除日志 | 🟡 待补完 | 上传时 `Signature` 仍为空，删除文件失败只留下 TODO（见 TODO #1/#2）。 |
| 更新检查 | `/api/check-update` | ✅ 已落地 | `updater.CheckUpdate` 使用语义化版本比较和校验和返回。 |
| 存储 | LocalStorage | ✅ 已落地 | file save/delete/url/exists 已实现并在 `main.go` 注入。 |
| 存储 | OSS / S3 适配器 | ⏳ 预留接口 | `storage.Storage` 提供接口，配置项存在但尚无实现分支。 |
| 中间件 | AdminMiddleware | ⚠️ 未接入 | 中间件已实现，但目前所有受保护路由仅做 JWT 校验，未强制管理员角色。 |

### 前端 Web 控制台

| 功能 | 状态 | 说明 |
|------|------|------|
| 登录页 | ✅ 可用 | `view/auth/login.vue` 已接入 `api.login` 与 Pinia 存储，并写入 localStorage。 |
| Token 刷新 | 🟡 未使用 | `api.refreshToken` 定义但尚未在 store/拦截器中调用，实际会在 Access Token 过期后直接失败。 |
| 仪表盘数据 | 🟡 静态占位 | `view/home/index.vue` 的统计卡片与最近动态均为硬编码示例，尚未调用后端。 |
| 导航/路由 | ⚠️ 不完整 | 路由仅注册 `/` 与 `/login`，但 `layout/default.vue` 的“应用管理”“系统设置”菜单会指向不存在的路由。 |
| 应用/版本/资产管理 UI | 🔲 未实现 | 目前没有对应的列表、表单或上传界面，API 只能通过 Postman/curl 使用。 |
| 注册/忘记密码 UI | 🔲 未实现 | 登录页仅展示提示文字，未提供表单或链接。 |

## ✅ 当前 TODO 列表

| # | 模块 | 文件 | 待办 | 建议思路 |
|---|------|------|------|----------|
| 1 | 资产上传 | `internal/asset/service.go` (行 133) | `Signature` 字段目前为空，需要实现文件签名生成与校验。 | 生成私钥签名或哈希签名存入字段，并在客户端校验；可考虑与 release key 绑定。 |
| 2 | 资产删除 | `internal/asset/service.go` (行 226) | 删除存储文件失败时仅注释提醒，需要记录日志。 | 注入 logger（`pkg/logger`）并输出错误详情，便于排查存储清理问题。 |
| 3 | 版本删除 | `internal/release/service.go` (行 137) | 删除 Release 时尚未同步删除关联 Assets。 | 在删除前查询并删除对应资产（含文件），或配置数据库外键级联，确保无孤立文件。 |

> 若新增 TODO，请同步更新此表，保持 README 成为 backlog 的单一事实来源。

## 🧪 测试与质量

- 后端可使用 `go test ./...` 运行单元测试。
- 前端使用 `pnpm test`（若添加测试工具）或 `pnpm lint` 保持规范。

## 🗺️ 后续展望

- 为资产签名、删除链路补全 TODO 后，可引入对象存储（S3、OSS 等）与 CI/CD 流水线。
- 可基于 `docs/UPDATE_CHECK.md` 实现自动更新客户端 SDK。
