# 部署指南 - AI热点追踪平台

## 架构概览

```
                    ┌─────────────┐
                    │   Client    │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │   Nginx     │  Port 80/443
                    │ (反向代理)   │
                    └──────┬──────┘
                           │
              ┌────────────┴────────────┐
              ↓                         ↓
    ┌─────────────────┐     ┌──────────────────┐
    │    Gateway      │     │  Content-Svc     │
    │   Port 8000     │     │  Port 8001       │
    │ - 认证/授权      │     │ - 文章管理        │
    │ - 用户管理       │     │ - 职业管理        │
    └────────┬────────┘     └────────┬─────────┘
             │                       │
             └───────────┬───────────┘
                         ↓
                  ┌─────────────┐
                  │   MySQL     │
                  │  Port 3306  │
                  └─────────────┘
```

## 三种部署方式

### 方式一: Docker Compose (推荐) ⭐

**优点**: 一键部署,环境隔离,易于扩展

```bash
# 1. 克隆项目
git clone <repository-url>
cd hot-ai-backend

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 文件,设置数据库密码等

# 3. 启动所有服务
docker-compose up -d

# 4. 查看服务状态
docker-compose ps

# 5. 查看日志
docker-compose logs -f
```

**访问地址**:
- API: http://localhost/api
- Health: http://localhost/health

**停止服务**:
```bash
docker-compose down
```

---

### 方式二: 使用启动脚本

**Windows**:
```powershell
# 1. 确保已安装 Go 和 Nginx
# 2. 配置数据库连接 (apps/services/content-svc/etc/content-svc.yaml)

# 3. 启动所有服务
.\start-with-nginx.bat

# 4. 停止所有服务
.\stop.bat
```

**Linux/Mac**:
```bash
# 1. 确保已安装 Go 和 Nginx
# 2. 赋予执行权限
chmod +x start-with-nginx.sh stop.sh

# 3. 启动所有服务
./start-with-nginx.sh

# 4. 停止所有服务
./stop.sh
```

---

### 方式三: 手动启动

**1. 启动 MySQL**
```bash
# 确保 MySQL 运行,并创建数据库
mysql -u root -p
CREATE DATABASE hot_ai CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

**2. 启动 Content Service**
```bash
cd hot-ai-backend
go run apps/services/content-svc/main.go -f apps/services/content-svc/etc/content-svc.yaml
```

**3. 启动 Gateway**
```bash
cd hot-ai-backend
go run apps/gateway/main.go -f apps/gateway/etc/gateway.yaml
```

**4. 启动 Nginx**
```bash
# 测试配置
nginx -t -c /path/to/nginx.conf

# 启动
nginx -c /path/to/nginx.conf
```

---

## 环境要求

### 必需组件

- **Go**: 1.22+
- **MySQL**: 8.0+
- **Nginx**: 1.20+ (可选,但强烈推荐)

### 可选组件

- **Redis**: 7.0+ (用于缓存)
- **Docker & Docker Compose**: 用于容器化部署

---

## 配置说明

### 1. Content Service 配置

文件: `apps/services/content-svc/etc/content-svc.yaml`

```yaml
Name: content-svc
Host: 0.0.0.0
Port: 8001

DataSource:
  MySQL:
    DSN: root:password@tcp(localhost:3306)/hot_ai?charset=utf8mb4&parseTime=True&loc=Local
```

### 2. Gateway 配置

文件: `apps/gateway/etc/gateway.yaml`

```yaml
Name: gateway
Host: 0.0.0.0
Port: 8000

Auth:
  AccessSecret: your-jwt-secret-key
  AccessExpire: 86400

Services:
  ContentSvc: http://localhost:8001  # Content Service 地址
```

### 3. Nginx 配置

文件: `nginx.conf`

主要配置上游服务器:
```nginx
upstream gateway_backend {
    server 127.0.0.1:8000;
}

upstream content_svc_backend {
    server 127.0.0.1:8001;
}
```

---

## 数据库初始化

### 方式一: Docker Compose (自动)

使用 Docker Compose 时,SQL 脚本会自动执行:

```bash
docker-compose up -d mysql
```

### 方式二: 手动执行 SQL 脚本

SQL 迁移脚本位于 `migrations/` 目录:

**1. 创建表结构**
```bash
mysql -u root -p < migrations/001_create_tables.sql
```

**2. 初始化文章数据**
```bash
mysql -u root -p < migrations/002_seed_article_data.sql
```

**3. 验证数据**
```sql
mysql -u root -p
USE hot_ai;
SHOW TABLES;
SELECT COUNT(*) FROM article_categories;  -- 应该返回 4
SELECT COUNT(*) FROM articles;             -- 应该返回 3
```

详细说明请查看: [migrations/README.md](migrations/README.md)

---

## 生产环境部署

### 1. 服务器要求

- **CPU**: 2核心+
- **内存**: 4GB+
- **磁盘**: 20GB+
- **操作系统**: Ubuntu 20.04 LTS / CentOS 8+

### 2. 安全配置

**防火墙设置**:
```bash
# Ubuntu
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# CentOS
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

**启用 HTTPS**:
```bash
# 使用 Let's Encrypt 获取免费证书
sudo apt-get install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### 3. 性能优化

**Nginx 优化** (`nginx.conf`):
```nginx
worker_processes auto;
worker_connections 2048;

gzip on;
gzip_comp_level 6;
gzip_types text/plain text/css application/json application/javascript;
```

**MySQL 优化** (`my.cnf`):
```ini
[mysqld]
max_connections = 200
innodb_buffer_pool_size = 1G
query_cache_size = 64M
```

### 4. 监控和日志

**日志轮转** (`/etc/logrotate.d/hot-ai`):
```
/var/log/nginx/*.log
/root/hot-ai-backend/logs/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0644 nginx nginx
    sharedscripts
    postrotate
        [ -f /var/run/nginx.pid ] && kill -USR1 `cat /var/run/nginx.pid`
    endscript
}
```

---

## 故障排查

### 常见问题

**1. 端口被占用**
```bash
# 查看端口占用
netstat -tlnp | grep :80
netstat -tlnp | grep :8000
netstat -tlnp | grep :8001

# 杀死占用进程
kill -9 <PID>
```

**2. 数据库连接失败**
```bash
# 测试数据库连接
mysql -h localhost -u root -p

# 检查 MySQL 状态
systemctl status mysql
```

**3. Nginx 502 错误**
```bash
# 检查后端服务是否运行
curl http://localhost:8000/health
curl http://localhost:8001/health

# 查看 Nginx 错误日志
tail -f /var/log/nginx/error.log
```

---

## 备份和恢复

### 数据库备份

```bash
# 备份
mysqldump -u root -p hot_ai > backup_$(date +%Y%m%d).sql

# 恢复
mysql -u root -p hot_ai < backup_20260406.sql
```

### 日志备份

```bash
# 压缩旧日志
tar -czf logs_backup_$(date +%Y%m%d).tar.gz logs/

# 清理旧日志
find logs/ -name "*.log" -mtime +30 -delete
```

---

## 更新部署

### Docker Compose 更新

```bash
# 1. 拉取最新代码
git pull

# 2. 重新构建并启动
docker-compose up -d --build

# 3. 查看日志确认启动成功
docker-compose logs -f
```

### 手动更新

```bash
# 1. 停止服务
./stop.sh

# 2. 拉取最新代码
git pull

# 3. 重新编译
go build -o bin/gateway ./apps/gateway
go build -o bin/content-svc ./apps/services/content-svc

# 4. 启动服务
./start-with-nginx.sh
```

---

## 扩展部署

### 多实例部署

修改 `nginx.conf` 实现负载均衡:

```nginx
upstream gateway_backend {
    least_conn;
    server 127.0.0.1:8000 weight=5;
    server 127.0.0.1:8001 weight=3;
    server 127.0.0.1:8002 backup;
}
```

### Kubernetes 部署

查看 `k8s/` 目录下的配置文件(待添加)。

---

## 技术支持

- **API 文档**: 查看 `API_DOCUMENTATION.md`
- **微服务架构**: 查看 `MICROSERVICES.md`
- **Nginx 配置**: 查看 `NGINX_GUIDE.md`

---

## 检查清单

部署前请确认:

- [ ] MySQL 已安装并运行
- [ ] Go 1.22+ 已安装
- [ ] Nginx 已安装(推荐使用)
- [ ] 数据库连接配置正确
- [ ] JWT Secret 已修改为强密码
- [ ] 防火墙规则已配置
- [ ] SSL 证书已配置(生产环境)
- [ ] 日志轮转已配置
- [ ] 数据库备份策略已制定

祝部署顺利! 🚀
