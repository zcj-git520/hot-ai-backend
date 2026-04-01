# AI 热点追踪平台 - 后端服务

## 技术栈

- **框架**: Go Zero v1.7
- **语言**: Go 1.21+
- **通信**: gRPC + Protobuf
- **服务发现**: Consul
- **数据库**: MySQL 8.0 + MongoDB 7.0
- **缓存**: Redis 7.0
- **消息队列**: Redis Stream → Kafka (演进)

## 项目结构

```
hot-ai-backend/
├── apps/
│   ├── gateway/              # API 网关
│   │   ├── main.go           # 入口文件
│   │   └── etc/              # 配置文件
│   │
│   └── services/
│       ├── content-svc/      # 内容服务
│       │   ├── main.go
│       │   └── etc/
│       │
│       └── crawler-svc/      # 采集服务
│           ├── main.go
│           └── etc/
│
├── internal/
│   ├── config/               # 配置结构
│   ├── handler/              # HTTP/gRPC 处理器
│   ├── model/                # 数据模型
│   └── svc/                  # 服务上下文
│
├── pkg/
│   └── proto/                # Protobuf 生成的代码
│
├── scripts/
│   └── gen-proto.sh          # Protobuf 生成脚本
│
└── deploy/
    └── docker/               # Docker 部署配置
```

## 快速开始

### 环境要求

- Go >= 1.21
- protoc (Protobuf 编译器)
- Docker & Docker Compose

### 1. 安装依赖

```bash
go mod download
```

### 2. 生成 Protobuf 代码

```bash
./scripts/gen-proto.sh
```

### 3. 启动基础设施

```bash
docker-compose up -d
```

### 4. 启动服务

```bash
# 启动 API 网关
go run apps/gateway/main.go

# 启动内容服务
go run apps/services/content-svc/main.go

# 启动采集服务
go run apps/services/crawler-svc/main.go
```

## API 文档

API 文档地址：http://localhost:8888/api/docs

## 服务监控

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001
- **Consul**: http://localhost:8500

## 开发规范

参考 [Go 语言开发规范](https://github.com/golang/go/wiki/CodeReviewComments)

## 许可证

MIT License
