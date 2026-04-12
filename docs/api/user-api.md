# 用户相关 API 接口文档

## 概述

本文档描述了用户管理和认证相关的 API 接口。所有接口均基于 RESTful 风格设计，使用 JSON 格式进行数据交换。

**基础路径**: `/api/v1`  
**认证方式**: JWT Bearer Token

---

## 目录

1. [认证接口](#认证接口)
   - [发送注册验证码](#发送注册验证码)
   - [用户注册](#用户注册)
   - [用户登录](#用户登录)
   - [刷新 Token](#刷新-token)
   - [用户登出](#用户登出)

---

## 数据结构

### UserInfo - 用户基本信息

```json
{
  "id": "string",
  "email": "string",
  "nickname": "string",
  "avatar": "string",
  "roles": ["string"]
}
```

**字段说明**:
- `id`: 用户 ID
- `email`: 邮箱地址
- `nickname`: 用户昵称
- `avatar`: 头像 URL（可选）
- `roles`: 角色列表

---

## 认证接口

### 发送注册验证码

向指定邮箱发送注册验证码，用于验证邮箱所有权。

**重要说明**:
- 验证码为 6 位数字
- 有效期 5 分钟
- 一次性使用，验证后自动失效
- 限流规则：同一邮箱 1 分钟内最多发送 1 次，1 小时内最多发送 5 次

**接口信息**:
- **路径**: `POST /api/v1/auth/send-registration-code`
- **认证**: 不需要
- **Content-Type**: `application/json`

**请求参数**:

```json
{
  "email": "user@example.com"
}
```

**请求字段说明**:
| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| email | string | 是 | 邮箱地址 | user@example.com |

**成功响应** (200 OK):

```json
{
  "message": "验证码已发送，请检查邮箱",
  "email": "user@example.com",
  "expires_in": 300
}
```

**响应字段说明**:
| 字段 | 类型 | 说明 |
|------|------|------|
| message | string | 提示信息 |
| email | string | 接收验证码的邮箱 |
| expires_in | integer | 验证码有效期（秒），默认 300 秒（5 分钟） |

**失败响应**:

**400 Bad Request** - 邮箱格式错误或已注册:
```json
{
  "code": 400,
  "message": "该邮箱已注册",
  "data": null
}
```

**429 Too Many Requests** - 操作过于频繁:
```json
{
  "code": 429,
  "message": "操作过于频繁，请稍后再试",
  "data": null
}
```

---

### 用户注册

创建新用户账号。

**重要说明**:
1. 必须先调用 `/api/v1/auth/send-registration-code` 发送邮箱验证码
2. 验证码有效期为 5 分钟，且只能使用一次
3. 密码需要输入两次以确保一致性

**接口信息**:
- **路径**: `POST /api/v1/auth/register`
- **认证**: 不需要
- **Content-Type**: `application/json`

**请求参数**:

```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "password_confirm": "SecurePass123!",
  "nickname": "用户名",
  "verification_code": "123456"
}
```

**请求字段说明**:
| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| email | string | 是 | 邮箱地址 | user@example.com |
| password | string | 是 | 密码（8-20 位，包含大小写字母 + 数字 + 特殊字符） | SecurePass123! |
| password_confirm | string | 是 | 确认密码（必须与 password 一致） | SecurePass123! |
| nickname | string | 是 | 用户昵称（2-20 个字符） | 用户名 |
| verification_code | string | 是 | 邮箱验证码（6 位数字） | 123456 |

**成功响应** (200 OK):

```json
{
  "user_id": "user-123",
  "message": "注册成功，请登录"
}
```

**失败响应** (400 Bad Request):

```json
{
  "code": 400,
  "message": "错误信息",
  "data": null
}
```

---

### 用户登录

用户登录获取访问令牌。

**接口信息**:
- **路径**: `POST /api/v1/auth/login`
- **认证**: 不需要
- **Content-Type**: `application/json`

**请求参数**:

```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**请求字段说明**:
| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| email | string | 是 | 邮箱地址 | user@example.com |
| password | string | 是 | 密码 | SecurePass123! |

**成功响应** (200 OK):

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 86400,
  "user": {
    "id": "user-123",
    "email": "user@example.com",
    "nickname": "测试用户",
    "avatar": "",
    "roles": ["user"]
  }
}
```

**响应字段说明**:
| 字段 | 类型 | 说明 |
|------|------|------|
| access_token | string | 访问令牌，用于后续请求认证 |
| refresh_token | string | 刷新令牌，用于获取新的 access_token |
| expires_in | integer | access_token 过期时间（秒），默认 86400 秒（24 小时） |
| user | object | 用户基本信息 |

**失败响应**:

**400 Bad Request** - 请求参数错误:
```json
{
  "code": 400,
  "message": "邮箱或密码格式错误",
  "data": null
}
```

**401 Unauthorized** - 认证失败:
```json
{
  "code": 401,
  "message": "邮箱或密码错误",
  "data": null
}
```

---

### 刷新 Token

使用 Refresh Token 获取新的 Access Token。

**接口信息**:
- **路径**: `POST /api/v1/auth/refresh`
- **认证**: 不需要
- **Content-Type**: `application/json`

**请求参数**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**请求字段说明**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| refresh_token | string | 是 | 之前获取的 refresh_token |

**成功响应** (200 OK):

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 86400
}
```

**响应字段说明**:
| 字段 | 类型 | 说明 |
|------|------|------|
| access_token | string | 新的访问令牌 |
| expires_in | integer | 过期时间（秒） |

**失败响应** (401 Unauthorized):

```json
{
  "code": 401,
  "message": "Refresh Token 无效或已过期",
  "data": null
}
```

---

### 用户登出

退出登录，使当前 Token 失效。

**接口信息**:
- **路径**: `POST /api/v1/auth/logout`
- **认证**: **需要** (Bearer Token)
- **Content-Type**: `application/json`

**请求头**:
```
Authorization: Bearer <access_token>
```

**成功响应** (200 OK):

```json
{
  "code": 200,
  "message": "登出成功",
  "data": {
    "success": true
  }
}
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权或认证失败 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 认证说明

除了注册和登录接口外，其他接口都需要在请求头中携带 JWT Token：

```
Authorization: Bearer <your_access_token>
```

Token 有效期为 24 小时，过期后可以使用 Refresh Token 刷新。

---

## 使用示例

### 1. 注册并登录

```bash
# 注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!",
    "nickname": "新用户"
  }'

# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

### 2. 使用 Token 访问受保护接口

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 3. 刷新 Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

---

## 更新日志

- **v1.0.0** (2026-03-31) - 初始版本
  - 用户注册功能
  - 用户登录功能
  - Token 刷新机制
  - 用户登出功能
