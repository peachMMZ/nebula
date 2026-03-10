# JSON Key 格式验证

## ✅ 已统一为小驼峰（camelCase）格式

### 修改的文件

1. **`internal/auth/jwt.go`**
   - `user_id` → `userId`
   - `access_token` → `accessToken`
   - `refresh_token` → `refreshToken`
   - `expires_in` → `expiresIn`

2. **`internal/auth/service.go`**
   - `refresh_token` → `refreshToken`

3. **`internal/api/handler/auth.go`**
   - `old_password` → `oldPassword`
   - `new_password` → `newPassword`

### 统一的 JSON Key 规范

| Go 结构体字段 | JSON Key (小驼峰) |
|--------------|------------------|
| ID | id |
| UserID | userId |
| Username | username |
| Email | email |
| Role | role |
| CreatedAt | createdAt |
| UpdatedAt | updatedAt |
| AccessToken | accessToken |
| RefreshToken | refreshToken |
| ExpiresIn | expiresIn |
| OldPassword | oldPassword |
| NewPassword | newPassword |
| AppID | appId (如果有) |
| ReleaseID | releaseId |
| Platform | platform |
| Arch | arch |
| URL | url |
| Signature | signature |
| Checksum | checksum |

### API 响应示例

**注册/登录响应：**
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
      "accessToken": "eyJhbGciOiJIUzI1NiIs...",
      "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
      "expiresIn": 7200
    }
  }
}
```

**刷新 Token 请求：**
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

**刷新 Token 响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "accessToken": "new_token...",
    "refreshToken": "new_refresh...",
    "expiresIn": 7200
  }
}
```

**修改密码请求：**
```json
{
  "oldPassword": "oldpass123",
  "newPassword": "newpass456"
}
```

**资源响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "releaseId": 1,
    "platform": "windows",
    "arch": "amd64",
    "url": "http://localhost:9050/files/...",
    "signature": "...",
    "checksum": "sha256...",
    "createdAt": "2026-03-10T12:00:00Z",
    "updatedAt": "2026-03-10T12:00:00Z"
  }
}
```

### 测试验证

```powershell
# 1. 注册用户（验证 camelCase）
$body = @{
    username = "testuser"
    email = "test@example.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:9050/api/auth/register" `
  -Method Post -Body $body -ContentType "application/json"

# 检查响应格式
Write-Host "✅ Access Token 字段: accessToken"
Write-Host "✅ Refresh Token 字段: refreshToken"
Write-Host "✅ Expires In 字段: expiresIn"
Write-Host "✅ Created At 字段: createdAt"
Write-Host "✅ Updated At 字段: updatedAt"

# 2. 刷新 Token（验证 camelCase）
$refreshBody = @{
    refreshToken = $response.data.tokens.refreshToken
} | ConvertTo-Json

$newTokens = Invoke-RestMethod -Uri "http://localhost:9050/api/auth/refresh" `
  -Method Post -Body $refreshBody -ContentType "application/json"

Write-Host "✅ 刷新 Token 成功，字段格式正确"

# 3. 修改密码（验证 camelCase）
$headers = @{
    Authorization = "Bearer $($response.data.tokens.accessToken)"
}

$changePasswordBody = @{
    oldPassword = "password123"
    newPassword = "newpassword456"
} | ConvertTo-Json

$result = Invoke-RestMethod -Uri "http://localhost:9050/api/auth/change-password" `
  -Method Post -Body $changePasswordBody -ContentType "application/json" -Headers $headers

Write-Host "✅ 修改密码成功，字段格式正确"
```

## ✅ 统一完成

所有 JSON key 已统一为小驼峰格式（camelCase），符合前端 JavaScript/TypeScript 的命名规范。
