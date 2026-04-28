# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AI 热点追踪平台 (AI Hot AI Backend) - 面向 AI 时代的职业导航系统后端，基于 Go Zero v1.7 微服务架构。

## Build & Run Commands

```bash
# Download dependencies
go mod download

# Generate Protobuf code (if modified)
./scripts/gen-proto.sh

# Run a specific service
go run apps/auth-sec/main.go
go run apps/content-svc/main.go
go run apps/crawler-svc/main.go
go run apps/learning-path-svc/main.go
go run apps/profession-svc/main.go
go run apps/tool-svc/main.go
```

## Architecture

### Microservices (apps/)

| Service | Port | Purpose |
|---------|------|---------|
| gateway (auth-sec) | 8000 | API 网关、认证 |
| content-svc | 8001 | 文章内容管理 |
| profession-svc | 8002 | 职业分析 |
| learning-path-svc | 8003 | 学习路径 |
| tool-svc | 8004 | 工具库 |
| crawler-svc | - | 爬虫采集（通过 Consul 服务发现） |

### Go Zero Service Structure

Each microservice follows Go Zero conventions:
```
apps/<service>/
├── main.go           # Entry point
├── etc/<service>.yaml # Configuration
├── handler/          # HTTP/gRPC handlers
├── logic/            # Business logic
├── model/            # Data models
├── repository/       # Data access layer
├── service/          # Service layer
└── internal/         # Package imports
```

### Key Dependencies

- **Database**: MySQL (GORM) + MongoDB
- **Cache**: Redis
- **Service Discovery**: Consul
- **Auth**: JWT (golang-jwt/jwt/v5)
- **Framework**: Go Zero v1.7

### Database Design (Tool Library Module)

Tool library uses MySQL with tables:
- `tools` - AI 工具信息
- `tool_categories` - 工具类别
- `tool_reviews` - 用户评测
- `prompt_templates` - 提示词模板
- `user_favorites` - 用户收藏
- `badges` / `user_badges` - 徽章系统

## Configuration

All services use YAML config files in `etc/` directory. Key config options:
- MySQL/Redis/Consul connection settings
- Service ports (8000-8004)
- JWT secret and expiration
- Log mode (console/file) and level

## Documentation

| 文档 | 路径 |
|------|------|
| 架构设计 | `docs/AI 热点追踪平台 - 软件架构设计文档.md` |
| 需求文档 | `docs/需求文档/AI 热点追踪平台 - 需求文档.md` |
| API 文档 | `docs/api/` |
| 数据库设计 | `docs/需求文档/DATABASE-DESIGN-TOOLS.md` |
| 爬虫配置 | `docs/爬虫配置指南.md` |
