# 微服务架构说明

## 架构概览

```
┌─────────────┐
│   Frontend   │  Nuxt.js (Port 3000)
└──────┬──────┘
       │ HTTP
       ↓
┌─────────────────────────────────────┐
│      API Gateway (Port 8000)        │
│  - 认证/授权                         │
│  - 用户管理                          │
│  - 反向代理到微服务                   │
└──┬──────────────┬───────────────────┘
   │              │
   │ HTTP         │ HTTP
   ↓              ↓
┌──────────┐  ┌──────────────────┐
│ User Svc │  │  Content Svc     │
│(内置)     │  │  (Port 8001)     │
│          │  │  - 文章管理       │
│          │  │  - 职业管理       │
│          │  │  - 分类管理       │
└──────────┘  └────────┬─────────┘
                       │
                  ┌────┴────┐
                  │  MySQL  │
                  └─────────┘
```

## 服务列表

### 1. API Gateway (`apps/gateway`)
**端口**: 8000

**职责**:
- 统一API入口
- JWT认证和授权
- 用户管理(注册/登录/个人资料)
- 反向代理到后端微服务
- 请求限流和日志记录

**路由**:
- `/api/auth/*` - 认证相关(本地处理)
- `/api/user/*` - 用户管理(本地处理)
- `/api/articles/*` - 文章相关(代理到 content-svc)
- `/api/professions/*` - 职业相关(代理到 content-svc)

### 2. Content Service (`apps/services/content-svc`)
**端口**: 8001

**职责**:
- 文章CRUD操作
- 职业信息管理
- 分类管理
- 搜索功能
- 阅读量统计

**API接口**:
- `GET /api/articles` - 获取文章列表
- `GET /api/articles/:id` - 获取文章详情
- `GET /api/articles/categories` - 获取文章分类
- `GET /api/professions` - 获取职业列表
- `GET /api/professions/:slug` - 获取职业详情
- `GET /api/professions/search?q=xxx` - 搜索职业
- `GET /api/professions/risk-levels` - 获取风险等级

## 启动顺序

1. **启动MySQL数据库**
2. **启动Content Service**
   ```bash
   go run apps/services/content-svc/main.go -f apps/services/content-svc/etc/content-svc.yaml
   ```
3. **启动Gateway**
   ```bash
   go run apps/gateway/main.go -f apps/gateway/etc/gateway.yaml
   ```

## 配置说明

### Content Service 配置
文件: `apps/services/content-svc/etc/content-svc.yaml`

```yaml
Name: content-svc
Host: 0.0.0.0
Port: 8001

DataSource:
  MySQL:
    DSN: root:password@tcp(localhost:3306)/hot_ai?charset=utf8mb4&parseTime=True&loc=Local
```

### Gateway 配置
文件: `apps/gateway/etc/gateway.yaml`

```yaml
Name: gateway
Host: 0.0.0.0
Port: 8000

Services:
  ContentSvc: http://localhost:8001  # Content Service地址
```

## 数据流示例

### 获取文章列表示例

```
1. 前端发起请求
   GET http://localhost:8000/api/articles?page=1&page_size=10

2. Gateway接收请求
   - 检查是否需要认证(此接口公开,无需认证)
   - 通过反向代理转发到 Content Service

3. Content Service处理
   - ArticleHandler.GetArticles() 接收请求
   - ArticleService.GetArticles() 执行业务逻辑
   - ArticleRepository.GetList() 查询数据库
   - 返回JSON响应

4. Gateway转发响应
   - 将Content Service的响应返回给前端

5. 前端接收响应
   {
     "code": 200,
     "message": "success",
     "data": {
       "articles": [...],
       "total": 100,
       "page": 1,
       "page_size": 10
     }
   }
```

## 扩展新的微服务

如果需要添加新的微服务(如学习路径服务):

1. **创建新服务目录**
   ```
   apps/services/learning-svc/
   ├── main.go
   ├── handler/
   ├── service/
   ├── repository/
   └── etc/learning-svc.yaml
   ```

2. **在Gateway中添加配置**
   ```yaml
   Services:
     ContentSvc: http://localhost:8001
     LearningSvc: http://localhost:8002  # 新增
   ```

3. **在Gateway中添加路由**
   ```go
   server.AddRoute(rest.Route{
       Method:  http.MethodGet,
       Path:    "/api/learning-paths/*",
       Handler: proxyHandler.ProxyToLearningSvc,
   })
   ```

## 优势

1. **独立部署**: 每个服务可以独立部署和扩展
2. **技术隔离**: 不同服务可以使用不同的技术栈
3. **故障隔离**: 单个服务故障不影响其他服务
4. **团队协作**: 不同团队可以并行开发不同服务
5. **易于维护**: 代码库更小,职责更清晰

## 注意事项

1. **服务依赖**: Gateway依赖Content Service,启动时需先启动Content Service
2. **网络通信**: 服务间通过HTTP通信,需确保网络畅通
3. **数据一致性**: 所有服务共享同一数据库,需注意事务处理
4. **配置管理**: 每个服务有自己的配置文件,需分别维护
