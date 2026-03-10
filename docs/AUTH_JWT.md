# JWT 认证系统文档

## 🔐 认证架构

Nebula 项目使用 **JWT (JSON Web Token)** 认证系统，提供完整的用户认证和授权功能。

### 核心特性

- ✅ **JWT 双 Token 机制**（Access Token + Refresh Token）
- ✅ **密码加密存储**（bcrypt）
- ✅ **灵活的 Token 过期时间配置**
- ✅ **中间件保护路由**
- ✅ **角色权限支持**（user/admin）
- ✅ **用户注册和登录**
- ✅ **修改密码**
- ✅ **Token 刷新**

## 📝 API 端点

### 1. 用户注册

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user": {
      "id": "uuid",
      "username": "testuser",
      "email": "test@example.com",
      "role": "user",
      "createdAt": "2026-03-10T12:00:00Z",
      "updatedAt": "2026-03-10T12:00:00Z"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIs...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
      "expires_in": 7200
    }
  }
}
```

### 2. 用户登录

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

**响应：** 同注册响应

### 3. 刷新 Token

```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "new_access_token",
    "refresh_token": "new_refresh_token",
    "expires_in": 7200
  }
}
```

### 4. 获取当前用户信息

```http
GET /api/auth/profile
Authorization: Bearer <access_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "uuid",
    "username": "testuser",
    "email": "test@example.com",
    "role": "user",
    "createdAt": "2026-03-10T12:00:00Z",
    "updatedAt": "2026-03-10T12:00:00Z"
  }
}
```

### 5. 修改密码

```http
POST /api/auth/change-password
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "old_password": "password123",
  "new_password": "newpassword456"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "password changed successfully"
}
```

## 🧪 PowerShell 测试示例

### 完整测试流程

```powershell
# 1. 注册用户
$registerBody = @{
    username = "testuser"
    email = "test@example.com"
    password = "password123"
} | ConvertTo-Json

$registerResponse = Invoke-RestMethod -Uri "http://localhost:9050/api/auth/register" `
  -Method Post -Body $registerBody -ContentType "application/json"

Write-Host "✅ 注册成功！"
Write-Host "用户 ID: $($registerResponse.data.user.id)"
Write-Host "Access Token: $($registerResponse.data.tokens.access_token.Substring(0, 20))..."

# 保存 token
$accessToken = $registerResponse.data.tokens.access_token
$refreshToken = $registerResponse.data.tokens.refresh_token

# 2. 登录（测试）
$loginBody = @{
    username = "testuser"
    password = "password123"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod -Uri "http://localhost:9050/api/auth/login" `
  -Method Post -Body $loginBody -ContentType "application/json"

Write-Host "✅ 登录成功！"

# 3. 获取用户信息
$headers = @{
    Authorization = "Bearer $accessToken"
}

$profile = Invoke-RestMethod -Uri "http://localhost:9050/api/auth/profile" `
  -Method Get -Headers $headers

Write-Host "✅ 获取用户信息成功！"
Write-Host "用户名: $($profile.data.username)"
Write-Host "邮箱: $($profile.data.email)"

# 4. 创建应用（测试认证保护）
$appBody = @{
    id = "test-app"
    name = "测试应用"
    description = "需要认证才能创建"
} | ConvertTo-Json

$app = Invoke-RestMethod -Uri "http://localhost:9050/api/apps" `
  -Method Post -Body $appBody -ContentType "application/json" -Headers $headers

Write-Host "✅ 创建应用成功（认证有效）！"

# 5. 刷新 Token
Start-Sleep -Seconds 2
$refreshBody = @{
    refresh_token = $refreshToken
} | ConvertTo-Json

$newTokens = Invoke-RestMethod -Uri "http://localhost:9050/api/auth/refresh" `
  -Method Post -Body $refreshBody -ContentType "application/json"

Write-Host "✅ Token 刷新成功！"

# 6. 修改密码
$changePasswordBody = @{
    old_password = "password123"
    new_password = "newpassword456"
} | ConvertTo-Json

$changeResult = Invoke-RestMethod -Uri "http://localhost:9050/api/auth/change-password" `
  -Method Post -Body $changePasswordBody -ContentType "application/json" -Headers $headers

Write-Host "✅ 密码修改成功！"

# 7. 测试未认证访问（应该失败）
try {
    Invoke-RestMethod -Uri "http://localhost:9050/api/apps" -Method Get
} catch {
    Write-Host "❌ 未认证访问被拒绝（正确）"
}
```

## 🔧 配置

### 环境变量

```bash
# JWT 密钥（生产环境必须修改！）
JWT_SECRET=your-secret-key-change-in-production

# Access Token 过期时间（默认 2 小时）
JWT_ACCESS_TOKEN_DURATION=7200  # 秒
# 或者
JWT_ACCESS_TOKEN_DURATION=2h

# Refresh Token 过期时间（默认 7 天）
JWT_REFRESH_TOKEN_DURATION=604800  # 秒
# 或者
JWT_REFRESH_TOKEN_DURATION=168h
```

### 推荐配置

**开发环境：**
- Access Token: 2 小时
- Refresh Token: 7 天

**生产环境：**
- Access Token: 15 分钟 - 1 小时
- Refresh Token: 30 天
- 强密钥（256位以上）

## 🛡️ 安全特性

### 1. 密码加密

```go
// 使用 bcrypt 加密
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 验证密码
bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
```

### 2. Token 验证

- ✅ 签名验证
- ✅ 过期时间检查
- ✅ 算法验证（HMAC-SHA256）

### 3. 中间件保护

```go
// 需要认证
authRequired.Use(auth.JWTMiddleware(jwtService))

// 需要管理员权限
adminRequired.Use(auth.JWTMiddleware(jwtService), auth.AdminMiddleware())
```

## 📊 路由保护状态

### 公开路由（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/auth/register | 用户注册 |
| POST | /api/auth/login | 用户登录 |
| POST | /api/auth/refresh | 刷新 Token |
| GET | /api/check-update | 检查更新 |

### 需要认证的路由

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/auth/profile | 获取用户信息 |
| POST | /api/auth/change-password | 修改密码 |
| GET/POST/PUT/DELETE | /api/apps/* | 应用管理 |
| GET/POST/PUT/DELETE | /api/releases/* | 版本管理 |
| GET/POST/PUT/DELETE | /api/assets/* | 资源管理 |

## 🔄 Token 刷新流程

```
Client                          Server
  |                               |
  |-------- Request API --------->| (Access Token)
  |<----- 401 Token Expired ------| ❌
  |                               |
  |--- Refresh Token Request ---->| (Refresh Token)
  |<--- New Access Token ---------|  ✅
  |                               |
  |-------- Request API --------->| (New Access Token)
  |<-------- Success -------------|  ✅
```

## 💡 客户端集成示例

### JavaScript/TypeScript

```typescript
class AuthService {
  private accessToken: string | null = null;
  private refreshToken: string | null = null;

  async login(username: string, password: string) {
    const response = await fetch('http://localhost:9050/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    });
    
    const data = await response.json();
    if (data.code === 0) {
      this.accessToken = data.data.tokens.access_token;
      this.refreshToken = data.data.tokens.refresh_token;
      
      // 保存到 localStorage
      localStorage.setItem('access_token', this.accessToken);
      localStorage.setItem('refresh_token', this.refreshToken);
    }
    return data;
  }

  async request(url: string, options: RequestInit = {}) {
    options.headers = {
      ...options.headers,
      'Authorization': `Bearer ${this.accessToken}`
    };

    let response = await fetch(url, options);
    
    // Token 过期，尝试刷新
    if (response.status === 401) {
      await this.refresh();
      options.headers['Authorization'] = `Bearer ${this.accessToken}`;
      response = await fetch(url, options);
    }
    
    return response.json();
  }

  async refresh() {
    const response = await fetch('http://localhost:9050/api/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: this.refreshToken })
    });
    
    const data = await response.json();
    if (data.code === 0) {
      this.accessToken = data.data.access_token;
      this.refreshToken = data.data.refresh_token;
      localStorage.setItem('access_token', this.accessToken);
      localStorage.setItem('refresh_token', this.refreshToken);
    }
  }
}
```

## 🚀 扩展功能

### 未来可添加

1. **Token 黑名单**：退出登录时将 token 加入黑名单
2. **双因素认证（2FA）**：增强安全性
3. **OAuth2 第三方登录**：支持 GitHub、Google 等
4. **权限细粒度控制**：基于资源的权限管理
5. **登录日志**：记录登录历史和异常
6. **IP 白名单**：限制访问来源
7. **账号锁定**：多次登录失败后锁定

## ✅ 测试验证

1. ✅ 注册新用户
2. ✅ 登录获取 token
3. ✅ 使用 token 访问受保护的 API
4. ✅ 刷新 token
5. ✅ 修改密码
6. ✅ 未认证访问被拒绝

完整的 JWT 认证系统已经实现并运行正常！🎉
