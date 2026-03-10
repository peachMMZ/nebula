# Asset API 使用示例

服务器地址: http://localhost:9050

## 📦 存储架构

### 可扩展的存储接口设计

项目采用接口化设计，支持多种存储后端：

```go
type Storage interface {
    Save(filename string, content io.Reader) (string, error)
    Delete(path string) error
    GetURL(path string) string
    Exists(path string) bool
}
```

**当前实现：**
- ✅ **LocalStorage** - 本地文件系统存储
- 🔲 **OSSStorage** - 阿里云 OSS（预留）
- 🔲 **S3Storage** - AWS S3（预留）

**切换存储方式只需：**
1. 实现 `Storage` 接口
2. 在配置中修改 `STORAGE_TYPE`
3. 无需修改业务代码

### 配置

通过环境变量配置存储：

```bash
# 本地存储（默认）
STORAGE_TYPE=local
STORAGE_BASE_PATH=./uploads
STORAGE_BASE_URL=http://localhost:9050/files

# 未来支持
# STORAGE_TYPE=oss
# STORAGE_TYPE=s3
```

## 📝 API 端点

### 1. 上传文件并创建资源

**重要：** 这是最常用的方式，上传文件并自动创建记录。

```bash
# 使用 curl
curl -X POST http://localhost:9050/api/releases/1/assets/upload \
  -F "file=@path/to/app.exe" \
  -F "platform=windows" \
  -F "arch=amd64"

# PowerShell 示例
$file = "D:\app.exe"
$uri = "http://localhost:9050/api/releases/1/assets/upload"

$form = @{
    file = Get-Item -Path $file
    platform = "windows"
    arch = "amd64"
}

Invoke-RestMethod -Uri $uri -Method Post -Form $form
```

**支持的文件类型：**
- Windows: `.exe`, `.msi`, `.zip`
- macOS: `.dmg`, `.pkg`, `.app`
- Linux: `.deb`, `.rpm`, `.appimage`, `.tar.gz`

**响应：**
```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 1,
    "releaseId": 1,
    "platform": "windows",
    "arch": "amd64",
    "url": "http://localhost:9050/files/releases/1/windows-amd64/app.exe",
    "checksum": "sha256哈希值",
    "createdAt": "2026-03-10T12:00:00Z"
  }
}
```

### 2. 获取所有资源列表

```bash
curl http://localhost:9050/api/assets

# PowerShell
Invoke-RestMethod -Uri "http://localhost:9050/api/assets"
```

### 3. 获取指定版本的资源列表

```bash
curl http://localhost:9050/api/releases/1/assets

# PowerShell
Invoke-RestMethod -Uri "http://localhost:9050/api/releases/1/assets"
```

### 4. 获取单个资源详情

```bash
curl http://localhost:9050/api/assets/1

# PowerShell
Invoke-RestMethod -Uri "http://localhost:9050/api/assets/1"
```

### 5. 创建资源记录（外部 URL）

如果文件托管在其他地方（如 GitHub Releases、CDN），可以只创建记录：

```bash
curl -X POST http://localhost:9050/api/assets \
  -H "Content-Type: application/json" \
  -d '{
    "releaseId": 1,
    "platform": "windows",
    "arch": "amd64",
    "url": "https://github.com/user/repo/releases/download/v1.0.0/app.exe",
    "checksum": "sha256...",
    "signature": "签名数据"
  }'

# PowerShell
$body = @{
    releaseId = 1
    platform = "windows"
    arch = "amd64"
    url = "https://github.com/user/repo/releases/download/v1.0.0/app.exe"
    checksum = "sha256..."
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:9050/api/assets" `
  -Method Post -Body $body -ContentType "application/json"
```

### 6. 更新资源信息

```bash
curl -X PUT http://localhost:9050/api/assets/1 \
  -H "Content-Type: application/json" \
  -d '{
    "platform": "windows",
    "arch": "arm64",
    "signature": "新签名"
  }'

# PowerShell
$body = @{
    platform = "windows"
    arch = "arm64"
    signature = "新签名"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:9050/api/assets/1" `
  -Method Put -Body $body -ContentType "application/json"
```

### 7. 删除资源

**注意：** 会同时删除文件和数据库记录。

```bash
curl -X DELETE http://localhost:9050/api/assets/1

# PowerShell
Invoke-RestMethod -Uri "http://localhost:9050/api/assets/1" -Method Delete
```

### 8. 下载资源文件

```bash
# 方式1：通过API端点下载
curl http://localhost:9050/api/assets/1/download -o app.exe

# 方式2：直接访问文件URL
curl http://localhost:9050/files/releases/1/windows-amd64/app.exe -o app.exe

# PowerShell - 方式1
Invoke-WebRequest -Uri "http://localhost:9050/api/assets/1/download" -OutFile "app.exe"

# PowerShell - 方式2（推荐，支持断点续传）
Invoke-WebRequest -Uri "http://localhost:9050/files/releases/1/windows-amd64/app.exe" -OutFile "app.exe"
```

## 🔄 完整工作流程

```powershell
# 1. 创建应用
$app = Invoke-RestMethod -Method POST -Uri "http://localhost:9050/api/apps" `
  -ContentType "application/json" `
  -Body '{"id":"my-app","name":"我的应用","description":"测试"}'

# 2. 创建版本
$release = Invoke-RestMethod -Method POST -Uri "http://localhost:9050/api/releases" `
  -ContentType "application/json" `
  -Body '{"appID":"my-app","version":"1.0.0","notes":"初始版本","channel":"stable","pubDate":"2026-03-10T10:00:00Z"}'

# 3. 上传 Windows 版本
$winForm = @{
    file = Get-Item -Path "D:\builds\app-windows.exe"
    platform = "windows"
    arch = "amd64"
}
$winAsset = Invoke-RestMethod -Uri "http://localhost:9050/api/releases/1/assets/upload" `
  -Method Post -Form $winForm

# 4. 上传 macOS 版本
$macForm = @{
    file = Get-Item -Path "D:\builds\app-macos.dmg"
    platform = "darwin"
    arch = "arm64"
}
$macAsset = Invoke-RestMethod -Uri "http://localhost:9050/api/releases/1/assets/upload" `
  -Method Post -Form $macForm

# 5. 上传 Linux 版本
$linuxForm = @{
    file = Get-Item -Path "D:\builds\app-linux"
    platform = "linux"
    arch = "amd64"
}
$linuxAsset = Invoke-RestMethod -Uri "http://localhost:9050/api/releases/1/assets/upload" `
  -Method Post -Form $linuxForm

# 6. 查看该版本的所有资源
$assets = Invoke-RestMethod -Uri "http://localhost:9050/api/releases/1/assets"
$assets.data | Format-Table id, platform, arch, url
```

## 📊 API 路由列表

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/assets | 获取所有资源列表 |
| GET | /api/assets/:id | 获取单个资源详情 |
| POST | /api/assets | 创建资源记录（外部URL） |
| PUT | /api/assets/:id | 更新资源信息 |
| DELETE | /api/assets/:id | 删除资源（含文件） |
| GET | /api/assets/:id/download | 下载资源文件 |
| GET | /api/releases/:id/assets | 获取版本的资源列表 |
| POST | /api/releases/:id/assets/upload | 上传文件并创建资源 |
| GET | /files/* | 静态文件访问 |

## 🔧 文件存储结构

```
uploads/
  └── releases/
      └── {releaseID}/
          ├── windows-amd64/
          │   └── app.exe
          ├── darwin-arm64/
          │   └── app.dmg
          └── linux-amd64/
              └── app
```

## 🔐 安全特性

1. **文件校验和**：自动计算 SHA256，确保文件完整性
2. **文件类型限制**：只允许特定扩展名的文件
3. **平台唯一性**：同一版本的同一平台+架构只能有一个资源
4. **级联删除**：删除资源时自动删除文件

## 🚀 扩展到云存储

未来要切换到 OSS/S3，只需：

1. **实现新的存储适配器**

```go
// internal/storage/oss.go
type OSSStorage struct {
    client *oss.Client
    bucket string
}

func (s *OSSStorage) Save(filename string, content io.Reader) (string, error) {
    // OSS 上传逻辑
}
// ... 实现其他接口方法
```

2. **在 main.go 中添加分支**

```go
switch cfg.Storage.Type {
case "local":
    stor, err = storage.NewLocalStorage(...)
case "oss":
    stor, err = storage.NewOSSStorage(...)
case "s3":
    stor, err = storage.NewS3Storage(...)
}
```

3. **业务代码无需修改** ✨

所有 Service 和 Handler 代码都基于接口编程，自动支持新的存储方式！
