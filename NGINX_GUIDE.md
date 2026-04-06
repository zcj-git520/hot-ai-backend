# Nginx 反向代理配置指南

## 架构说明

```
客户端请求
    ↓
┌──────────────┐
│   Nginx      │  Port 80 (HTTP) / 443 (HTTPS)
│  (反向代理)   │
└──────┬───────┘
       │
       ├──────────────────┐
       ↓                  ↓
┌─────────────┐   ┌──────────────┐
│   Gateway   │   │ Content-Svc  │
│  Port 8000  │   │  Port 8001   │
└─────────────┘   └──────────────┘
```

## 快速开始

### Windows

**1. 安装 Nginx**
```powershell
# 下载 Nginx for Windows
# http://nginx.org/en/download.html
# 解压到 C:\nginx

# 添加到系统PATH
# 系统属性 → 高级 → 环境变量 → Path → 添加 C:\nginx
```

**2. 启动所有服务**
```powershell
# 在项目根目录执行
.\start-with-nginx.bat
```

**3. 停止所有服务**
```powershell
.\stop.bat
```

### Linux/Mac

**1. 安装 Nginx**
```bash
# Ubuntu/Debian
sudo apt-get install nginx

# CentOS/RHEL
sudo yum install nginx

# MacOS
brew install nginx
```

**2. 启动所有服务**
```bash
chmod +x start-with-nginx.sh stop.sh
./start-with-nginx.sh
```

**3. 停止所有服务**
```bash
./stop.sh
```

### Docker Compose (推荐)

**一键启动所有服务(包括Nginx)**
```bash
docker-compose up -d
```

**查看日志**
```bash
docker-compose logs -f nginx
docker-compose logs -f gateway
docker-compose logs -f content-svc
```

**停止服务**
```bash
docker-compose down
```

## Nginx 配置说明

### 主要配置文件: `nginx.conf`

**关键配置项:**

1. **上游服务器定义**
```nginx
upstream gateway_backend {
    server 127.0.0.1:8000;
    keepalive 32;
}

upstream content_svc_backend {
    server 127.0.0.1:8001;
    keepalive 32;
}
```

2. **API代理配置**
```nginx
location /api/ {
    proxy_pass http://gateway_backend;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

3. **健康检查端点**
```nginx
location /health {
    return 200 '{"status":"ok"}';
    add_header Content-Type application/json;
}
```

## 常用命令

### Nginx 管理

```bash
# 测试配置文件
nginx -t

# 启动 Nginx
nginx

# 停止 Nginx
nginx -s stop

# 重启 Nginx
nginx -s reload

# 查看 Nginx 状态
ps aux | grep nginx

# 查看错误日志
tail -f /var/log/nginx/error.log  # Linux
tail -f logs/error.log            # Windows
```

### 服务管理

```bash
# 查看所有服务进程
ps aux | grep -E "nginx|gateway|content-svc"

# 查看端口占用
netstat -tlnp | grep -E "80|8000|8001"  # Linux
netstat -ano | findstr "80 8000 8001"    # Windows
```

## 访问地址

启动成功后,可以通过以下地址访问:

- **API入口**: http://localhost/api
- **健康检查**: http://localhost/health
- **Gateway直接访问**: http://localhost:8000/api
- **Content-Svc直接访问**: http://localhost:8001/api

## 生产环境配置

### 1. 启用 HTTPS

取消 `nginx.conf` 中 HTTPS 配置的注释:

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # SSL 优化配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    location /api/ {
        proxy_pass http://gateway_backend;
        # ... 其他配置
    }
}

# HTTP 自动跳转 HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

### 2. 负载均衡

修改 `nginx.conf`:

```nginx
upstream gateway_backend {
    least_conn;  # 最少连接算法
    server 127.0.0.1:8000 weight=5;
    server 127.0.0.1:8001 weight=3;
    server 127.0.0.1:8002 backup;
}
```

### 3. 限流保护

```nginx
http {
    # 定义限流区域
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

    server {
        location /api/ {
            limit_req zone=api_limit burst=20 nodelay;
            proxy_pass http://gateway_backend;
        }
    }
}
```

### 4. 缓存静态资源

```nginx
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
    expires 30d;
    add_header Cache-Control "public, immutable";
    proxy_pass http://gateway_backend;
}
```

## 故障排查

### 问题1: Nginx 启动失败

```bash
# 检查配置文件语法
nginx -t

# 查看错误日志
tail -f logs/error.log

# 检查端口是否被占用
netstat -tlnp | grep :80
```

### 问题2: 502 Bad Gateway

```bash
# 检查后端服务是否运行
curl http://localhost:8000/health
curl http://localhost:8001/health

# 查看 Nginx 错误日志
tail -f logs/error.log

# 检查防火墙设置
sudo ufw status  # Linux
```

### 问题3: 请求超时

调整 `nginx.conf` 中的超时设置:

```nginx
location /api/ {
    proxy_connect_timeout 120s;
    proxy_send_timeout 120s;
    proxy_read_timeout 120s;
}
```

## 性能优化

### 1. 开启 Gzip 压缩

已在配置文件中启用,可根据需要调整压缩级别。

### 2. 调整 Worker 进程

```nginx
worker_processes auto;  # 自动根据CPU核心数设置
worker_connections 2048;  # 增加并发连接数
```

### 3. 启用缓存

```nginx
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m max_size=1g inactive=60m;

location /api/articles {
    proxy_cache my_cache;
    proxy_cache_valid 200 10m;
    proxy_cache_use_stale error timeout updating;
}
```

## 监控和日志

### 日志位置

- **Access Log**: `logs/access.log`
- **Error Log**: `logs/error.log`

### 日志分析

```bash
# 查看访问量最高的IP
awk '{print $1}' logs/access.log | sort | uniq -c | sort -nr | head -10

# 查看404错误
grep " 404 " logs/access.log

# 查看今日访问量
grep $(date +%d/%b/%Y) logs/access.log | wc -l
```

## 安全建议

1. **隐藏 Nginx 版本号**
```nginx
server_tokens off;
```

2. **添加安全头**
```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
```

3. **限制请求体大小**
```nginx
client_max_body_size 10M;
```

4. **启用 HTTPS** (见上文)

## 更多资源

- [Nginx 官方文档](https://nginx.org/en/docs/)
- [Nginx 最佳实践](https://www.nginx.com/blog/nginx-best-practices/)
- [Let's Encrypt - 免费SSL证书](https://letsencrypt.org/)
