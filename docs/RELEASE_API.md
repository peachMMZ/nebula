# Release API 使用示例

服务器地址: http://localhost:9050

## 1. 创建应用（前置步骤）

```bash
curl -X POST http://localhost:9050/api/apps \
  -H "Content-Type: application/json" \
  -d '{
    "id": "my-app",
    "name": "我的应用",
    "description": "这是一个测试应用"
  }'
```

## 2. 创建版本发布

```bash
curl -X POST http://localhost:9050/api/releases \
  -H "Content-Type: application/json" \
  -d '{
    "appID": "my-app",
    "version": "1.0.0",
    "notes": "第一个正式版本\n- 新增功能A\n- 修复bug B",
    "channel": "stable",
    "pubDate": "2026-03-10T10:00:00Z"
  }'
```

## 3. 获取所有版本列表

```bash
curl http://localhost:9050/api/releases
```

## 4. 获取指定应用的版本列表

```bash
curl http://localhost:9050/api/apps/my-app/releases
```

## 5. 按渠道筛选版本

```bash
# 获取 stable 渠道的版本
curl "http://localhost:9050/api/apps/my-app/releases?channel=stable"
```

## 6. 获取最新版本

```bash
curl http://localhost:9050/api/apps/my-app/releases/latest

# 获取指定渠道的最新版本
curl "http://localhost:9050/api/apps/my-app/releases/latest?channel=stable"
```

## 7. 获取单个版本详情

```bash
curl http://localhost:9050/api/releases/1
```

## 8. 更新版本信息

```bash
curl -X PUT http://localhost:9050/api/releases/1 \
  -H "Content-Type: application/json" \
  -d '{
    "version": "1.0.1",
    "notes": "更新内容",
    "channel": "stable",
    "pubDate": "2026-03-10T12:00:00Z"
  }'
```

## 9. 删除版本

```bash
curl -X DELETE http://localhost:9050/api/releases/1
```

## 完整测试流程

```powershell
# 1. 创建应用
Invoke-RestMethod -Method POST -Uri "http://localhost:9050/api/apps" `
  -ContentType "application/json" `
  -Body '{"id":"test-app","name":"测试应用","description":"测试用"}'

# 2. 创建第一个版本
Invoke-RestMethod -Method POST -Uri "http://localhost:9050/api/releases" `
  -ContentType "application/json" `
  -Body '{"appID":"test-app","version":"1.0.0","notes":"初始版本","channel":"stable","pubDate":"2026-03-10T10:00:00Z"}'

# 3. 创建第二个版本
Invoke-RestMethod -Method POST -Uri "http://localhost:9050/api/releases" `
  -ContentType "application/json" `
  -Body '{"appID":"test-app","version":"1.1.0","notes":"功能更新","channel":"stable","pubDate":"2026-03-10T12:00:00Z"}'

# 4. 查看该应用的所有版本
Invoke-RestMethod -Uri "http://localhost:9050/api/apps/test-app/releases"

# 5. 获取最新版本
Invoke-RestMethod -Uri "http://localhost:9050/api/apps/test-app/releases/latest"
```

## API 路由列表

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/releases | 获取所有版本列表 |
| GET | /api/releases/:id | 获取单个版本详情 |
| POST | /api/releases | 创建新版本 |
| PUT | /api/releases/:id | 更新版本信息 |
| DELETE | /api/releases/:id | 删除版本 |
| GET | /api/apps/:id/releases | 获取指定应用的版本列表 |
| GET | /api/apps/:id/releases/latest | 获取应用的最新版本 |

## 响应格式

成功响应:
```json
{
  "code": 0,
  "msg": "ok",
  "data": { ... }
}
```

错误响应:
```json
{
  "code": 500,
  "msg": "错误信息"
}
```
