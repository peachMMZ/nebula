# 更新检查 API 文档

## 🔧 修复内容

### 已修复的 Bug

1. ✅ **版本比较逻辑错误**
   - 之前：使用字符串比较（`"1.10.0" < "1.9.0"` 错误）
   - 现在：实现语义化版本号比较（semver）

2. ✅ **Asset 查询错误处理缺失**
   - 之前：查询失败时返回空 URL
   - 现在：正确处理错误并返回明确的错误信息

3. ✅ **参数验证缺失**
   - 现在：验证所有必需参数（app, version, platform, arch）

4. ✅ **排序逻辑优化**
   - 之前：按 version 字符串排序
   - 现在：按 pub_date 排序，更准确

5. ✅ **响应格式统一**
   - 现在使用统一的响应格式，并添加 checksum 字段

## 📝 API 使用

### 端点

```
GET /api/check-update
```

### 请求参数

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| app | string | ✅ | 应用 ID |
| version | string | ✅ | 当前版本号 |
| platform | string | ✅ | 平台（windows/darwin/linux） |
| arch | string | ✅ | 架构（amd64/arm64/386） |

### 响应格式

**有更新：**
```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "update": true,
    "version": "1.2.0",
    "notes": "更新说明\n- 新功能\n- Bug修复",
    "url": "http://localhost:9050/files/releases/2/windows-amd64/app.exe",
    "checksum": "sha256哈希值"
  }
}
```

**无需更新：**
```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "update": false
  }
}
```

**错误响应：**
```json
{
  "code": 500,
  "msg": "no asset found for this platform and architecture"
}
```

## 🧪 测试示例

### PowerShell 测试

```powershell
# 测试更新检查（需要更新）
$params = @{
    app = "test-app"
    version = "1.0.0"
    platform = "windows"
    arch = "amd64"
}
$query = ($params.GetEnumerator() | ForEach-Object { "$($_.Key)=$($_.Value)" }) -join "&"
Invoke-RestMethod -Uri "http://localhost:9050/api/check-update?$query"

# 测试更新检查（已是最新）
$params = @{
    app = "test-app"
    version = "2.0.0"
    platform = "windows"
    arch = "amd64"
}
$query = ($params.GetEnumerator() | ForEach-Object { "$($_.Key)=$($_.Value)" }) -join "&"
Invoke-RestMethod -Uri "http://localhost:9050/api/check-update?$query"

# 测试不存在的平台
$params = @{
    app = "test-app"
    version = "1.0.0"
    platform = "windows"
    arch = "arm"  # 假设没有这个架构的版本
}
$query = ($params.GetEnumerator() | ForEach-Object { "$($_.Key)=$($_.Value)" }) -join "&"
Invoke-RestMethod -Uri "http://localhost:9050/api/check-update?$query"
```

### curl 测试

```bash
# 测试更新检查
curl "http://localhost:9050/api/check-update?app=test-app&version=1.0.0&platform=windows&arch=amd64"

# 测试参数缺失
curl "http://localhost:9050/api/check-update?app=test-app"
```

## 📦 版本号比较规则

实现了完整的语义化版本号（semver）比较：

### 支持的版本格式

- `1.0.0` - 标准格式
- `v1.0.0` - 带 v 前缀
- `1.0.0-beta.1` - 预发布版本
- `1.0.0+build.123` - 构建元数据

### 比较规则

1. **主版本号优先**：`2.0.0 > 1.9.9`
2. **次版本号次之**：`1.10.0 > 1.9.0`
3. **修订号最后**：`1.0.10 > 1.0.9`
4. **预发布版本**：`1.0.0 > 1.0.0-rc.1 > 1.0.0-beta.1 > 1.0.0-alpha.1`

### 示例

```go
// pkg/util/version.go

CompareVersion("1.10.0", "1.9.0")      // 返回 1 (1.10.0 > 1.9.0)
CompareVersion("2.0.0", "1.99.99")     // 返回 1 (2.0.0 > 1.99.99)
CompareVersion("1.0.0", "1.0.0-beta")  // 返回 1 (正式版 > 预发布版)
CompareVersion("v1.0.0", "1.0.0")      // 返回 0 (相等)

IsNewerVersion("1.0.0", "1.1.0")       // 返回 true
IsNewerVersion("1.1.0", "1.0.0")       // 返回 false
```

## 🔄 完整更新流程

### 1. 准备数据

```powershell
# 创建应用
Invoke-RestMethod -Method POST -Uri "http://localhost:9050/api/apps" `
  -ContentType "application/json" `
  -Body '{"id":"my-app","name":"我的应用","description":"测试"}'

# 创建版本 1.0.0
Invoke-RestMethod -Method POST -Uri "http://localhost:9050/api/releases" `
  -ContentType "application/json" `
  -Body '{"appID":"my-app","version":"1.0.0","notes":"初始版本","channel":"stable","pubDate":"2026-03-10T10:00:00Z"}'

# 上传 1.0.0 的 Windows 版本
$form = @{
    file = Get-Item -Path "D:\app-v1.0.0.exe"
    platform = "windows"
    arch = "amd64"
}
Invoke-RestMethod -Uri "http://localhost:9050/api/releases/1/assets/upload" -Method Post -Form $form

# 创建版本 1.1.0
Invoke-RestMethod -Method POST -Uri "http://localhost:9050/api/releases" `
  -ContentType "application/json" `
  -Body '{"appID":"my-app","version":"1.1.0","notes":"新增功能","channel":"stable","pubDate":"2026-03-11T10:00:00Z"}'

# 上传 1.1.0 的 Windows 版本
$form = @{
    file = Get-Item -Path "D:\app-v1.1.0.exe"
    platform = "windows"
    arch = "amd64"
}
Invoke-RestMethod -Uri "http://localhost:9050/api/releases/2/assets/upload" -Method Post -Form $form
```

### 2. 客户端检查更新

```powershell
# 用户当前使用 1.0.0，检查更新
$response = Invoke-RestMethod -Uri "http://localhost:9050/api/check-update?app=my-app&version=1.0.0&platform=windows&arch=amd64"

if ($response.data.update) {
    Write-Host "发现新版本: $($response.data.version)"
    Write-Host "更新说明: $($response.data.notes)"
    Write-Host "下载地址: $($response.data.url)"
    Write-Host "校验和: $($response.data.checksum)"
    
    # 下载更新
    Invoke-WebRequest -Uri $response.data.url -OutFile "app-update.exe"
    
    # 验证校验和（可选）
    $hash = (Get-FileHash -Path "app-update.exe" -Algorithm SHA256).Hash
    if ($hash.ToLower() -eq $response.data.checksum) {
        Write-Host "✅ 文件校验通过"
    } else {
        Write-Host "❌ 文件校验失败"
    }
} else {
    Write-Host "✅ 当前已是最新版本"
}
```

## 🛡️ 错误处理

所有可能的错误情况：

| 错误信息 | 原因 | 解决方法 |
|---------|------|---------|
| `app is required` | 缺少 app 参数 | 提供应用 ID |
| `version is required` | 缺少 version 参数 | 提供当前版本号 |
| `platform is required` | 缺少 platform 参数 | 提供平台信息 |
| `arch is required` | 缺少 arch 参数 | 提供架构信息 |
| `no release found for this app` | 应用没有发布版本 | 先创建版本发布 |
| `no asset found for this platform and architecture` | 没有对应平台的安装包 | 上传对应平台的资源 |
| `asset URL is empty` | 资源记录存在但 URL 为空 | 检查资源数据完整性 |

## ✅ 改进总结

### 修复前的问题

```go
// ❌ 错误的版本比较
if req.Version >= latest.Version {  // "1.10.0" < "1.9.0" ！
    return &CheckResponse{Update: false}, nil
}

// ❌ 没有错误处理
db.Where(...).First(&asset)  // 忽略了错误
return &CheckResponse{URL: asset.URL}  // URL 可能为空
```

### 修复后

```go
// ✅ 正确的版本比较
if !util.IsNewerVersion(req.Version, latest.Version) {
    return &CheckResponse{Update: false}, nil
}

// ✅ 完整的错误处理
err = db.Where(...).First(&ast).Error
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errors.New("no asset found...")
    }
    return nil, err
}
```

现在更新检查功能已经完全可靠！🎉
