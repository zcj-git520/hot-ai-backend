# 学习路径微服务 (Learning Path Service)

## 概述

学习路径微服务负责管理学习路径、章节和学习进度跟踪，为 AI 热点追踪平台用户提供结构化的学习体验。

## 功能特性

- ✅ 学习路径管理（CRUD、筛选、排序、推荐）
- ✅ 章节管理（内容管理、跳转导航）
- ✅ 学习进度跟踪（用户进度保存、完成章节统计）
- ✅ 难度等级管理（入门/进阶/高级）
- ✅ 学生统计数据（学习人数、完成率等）

## 项目结构

```
learning-path-svc/
├── etc/
│   └── learning-path-svc.yaml      # 配置文件
├── handler/
│   └── learning_path_handler.go    # HTTP 处理器
└── main.go                         # 服务入口

internal/
├── repository/
│   └── learning_path_repository.go # 数据访问层（已存在）
└── service/
    └── learning_path_service.go    # 业务逻辑层

models/
├── LearningPath                    # 学习路径模型
├── PathChapter                     # 章节模型
├── LearningProgress                # 学习进度模型
└── LearningPathManagement          # 管理数据模型
```

## API 接口

### 1. 学习路径管理

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/learning-paths` | GET | 获取路径列表（支持分页和难度筛选） |
| `/api/learning-paths/:id` | GET | 根据 ID 获取路径详情 |
| `/api/learning-paths/slug/:slug` | GET | 根据 slug 获取路径详情 |
| `/api/learning-paths/featured` | GET | 获取推荐路径 |
| `/api/learning-paths/levels` | GET | 获取难度等级信息 |

### 2. 章节管理

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/learning-paths/:path_id/chapters` | GET | 获取路径的所有章节 |
| `/api/learning-paths/chapters/:chapter_id` | GET | 根据 ID 获取章节详情 |
| `/api/learning-paths/:path_slug/:chapter_slug` | GET | 根据 slug 获取章节详情（含前后章） |

### 3. 学习进度管理

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/learning-paths/progress` | GET | 获取用户学习进度 |
| `/api/learning-paths/completed-chapters` | GET | 获取已完成章节列表 |
| `/api/learning-paths/progress` | POST | 保存学习进度 |
| `/api/learning-paths/dashboard` | GET | 获取学习仪表盘 |

## 快速开始

### 环境要求

- Go 1.22+
- MySQL 8.0+
- Docker & Docker Compose（可选）

### 安装依赖

```bash
# 安装 Go 依赖
go mod download
```

### 配置数据库

确保你的 `migrations/005_create_learning_paths.sql` 和 `migrations/006_seed_learning_path_data.sql` 已执行。

### 运行服务

#### 方式 1：本地运行

```bash
export MYSQL_DSN="hot_ai:hotai123@tcp(localhost:3306)/hot_ai?charset=utf8mb4&parseTime=True&loc=Local"

go run apps/services/learning-path-svc/main.go \
  -f apps/services/learning-path-svc/etc/learning-path-svc.yaml
```

#### 方式 2：使用启动脚本

```bash
# 本地模式
./run-learning-path-svc.sh -l

# Docker 模式（推荐）
./run-learning-path-svc.sh -d

# 重新构建并启动
./run-learning-path-svc.sh -d -b
```

#### 方式 3：Docker Compose

```bash
# 单独启动服务
docker-compose up -d learning-path-svc

# 重新构建并启动
docker-compose build learning-path-svc
docker-compose up -d learning-path-svc
```

### 验证服务

```bash
# 健康检查
curl http://localhost:8003/health

# 获取路径列表
curl http://localhost:8003/api/learning-paths

# 获取推荐路径
curl http://localhost:8003/api/learning-paths/featured
```

## 配置说明

### 配置文件 (etc/learning-path-svc.yaml)

```yaml
Name: learning-path-svc
Host: 0.0.0.0
Port: 8003

DataSource:
  MySQL:
    DSN: "${MYSQL_DSN}"  # 使用环境变量
```

### 环境变量

| 变量 | 说明 | 示例 |
|------|------|------|
| MYSQL_DSN | MySQL 连接字符串 | `user:pass@tcp(host:3306)/db?charset=utf8mb4&parseTime=True&loc=Local` |

## 数据库表结构

学习路径服务使用以下数据库表：

1. **learning_paths** - 学习路径主表
2. **path_chapters** - 章节表
3. **learning_progress** - 学习进度表
4. **learning_path_management** - 管理数据统计表

详见：`migrations/005_create_learning_paths.sql`

## 数据模型

### LearningPath

```go
type LearningPath struct {
    ID             uint      // 路径 ID
    Title          string    // 标题
    Slug           string    // URL 友好标识
    Icon           string    // Emoji 图标
    Description    string    // 描述
    Difficulty     string    // 难度等级
    LearningGoals  string    // 学习目标（JSON）
    TargetAudience string    // 适合人群（JSON）
    // ... 其他字段
    Chapters       []PathChapter  // 关联章节（非持久化）
}
```

### PathChapter

```go
type PathChapter struct {
    ID            uint      // 章节 ID
    PathID        uint      // 路径 ID
    Title         string    // 标题
    Slug          string    // URL 友好标识
    ContentType   string    // 内容类型（article/video/practice/external）
    Content       string    // 内容正文（Markdown）
    VideoURL      string    // 视频 URL
    ExternalLinks string    // 外部链接（JSON）
    // ... 其他字段
}
```

### LearningProgress

```go
type LearningProgress struct {
    ID          uint      // 进度 ID
    UserID      string    // 用户 ID
    SessionID   string    // 会话 ID
    PathID      uint      // 路径 ID
    ChapterID   uint      // 章节 ID
    Status      string    // 状态（in_progress/completed）
    TimeSpent   int       // 学习时长（分钟）
    Notes       string    // 学习笔记
    // ... 其他字段
}
```

## 完整的 API 文档

查看详细的接口文档：[LEARNING_PATH_API.md](../../LEARNING_PATH_API.md)

## 开发指南

### 添加新接口

1. 在 `handler/learning_path_handler.go` 中添加处理器方法
2. 在 `main.go` 的 `registerRoutes` 中注册路由
3. （如果需要）在 `service/learning_path_service.go` 中添加业务逻辑

### 添加新功能

1. 先在 `internal/repository/learning_path_repository.go` 中实现数据访问
2. 在 `internal/service/learning_path_service.go` 中实现业务逻辑
3. 在 `handler/learning_path_handler.go` 中暴露 HTTP 接口

## 监控与日志

### 健康检查

```bash
# 健康检查接口
curl http://localhost:8003/health
```

### 查看日志

```bash
# Docker 容器日志
docker logs -f hot-ai-learning-path-svc

# 本地日志（如果有配置）
tail -f logs/learning-path-svc.log
```

## 性能优化

1. **缓存策略**：路径列表和详情可缓存 5-10 分钟
2. **数据库索引**：确保 `learning_paths.slug`、`path_chapters.path_id` 等字段有索引
3. **连接池**：使用 MySQL 连接池，避免频繁创建连接

## 故障排查

### 服务无法启动

1. 检查数据库连接配置
2. 验证表结构是否已创建：`migrations/005_create_learning_paths.sql`
3. 查看端口 8003 是否被占用

### API 返回 404

1. 确认请求 URL 正确
2. 检查路由是否已注册（main.go）
3. 查看路径是否存在（slug 是否正确）

### 数据库错误

1. 检查 MySQL 服务是否运行
2. 验证数据库迁移是否成功
3. 查看应用日志获取详细错误信息

## 集成测试

```bash
# 测试所有接口
cd scripts
go run test_learning_path_api.go

# 或使用 Postman 导入集合
# Postman 集合文件: docs/learning-path-api.postman_collection.json
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -f Dockerfile.learning-path-svc -t hot-ai/learning-path-svc .

# 运行容器
docker run -d \
  --name learning-path-svc \
  -e MYSQL_DSN="user:pass@tcp(db:3306)/hot_ai?charset=utf8mb4&parseTime=True&loc=Local" \
  -p 8003:8003 \
  hot-ai/learning-path-svc
```

### Kubernetes 部署

参考: `k8s/learning-path-svc-deployment.yaml`

## 相关文档

- [完整API文档](../../LEARNING_PATH_API.md)
- [数据库迁移](../../migrations/005_create_learning_paths.sql)
- [Docker 部署指南](../DEPLOYMENT.md)
- [Nginx 配置](../../nginx.conf)

## 支持

遇到问题？请查看：[故障排查](#故障排查) 或提交 Issue

## 贡献

欢迎提交 Pull Request 或 Issue！
