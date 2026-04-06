# 数据库迁移脚本

## 📁 文件说明

```
migrations/
├── 001_create_tables.sql       # 创建所有表结构
├── 002_seed_article_data.sql   # 文章相关初始数据
└── README.md                   # 本文档
```

## 🚀 使用方法

### 方式一: 按顺序执行所有脚本

```bash
# 1. 创建表结构
mysql -u root -p < migrations/001_create_tables.sql

# 2. 初始化文章数据
mysql -u root -p < migrations/002_seed_article_data.sql
```

### 方式二: Docker Compose (自动)

Docker Compose 会自动执行 `migrations/` 目录下的脚本(需配置卷映射):

```yaml
volumes:
  - ./migrations:/docker-entrypoint-initdb.d
```

## 📊 脚本说明

### 001_create_tables.sql

创建所有必需的表:

**用户相关:**
- users - 用户表
- roles - 角色表
- permissions - 权限表
- user_roles - 用户角色关联
- role_permissions - 角色权限关联
- refresh_tokens - 刷新Token
- favorites - 收藏表

**内容相关:**
- article_categories - 文章分类
- articles - 文章表
- professions - 职业表
- risk_level_info - 风险等级

**爬虫相关:**
- crawler_sources - 爬虫源配置
- crawler_fetch_logs - 爬取日志

### 002_seed_article_data.sql

初始化文章相关数据:

**文章分类 (4条):**
- AI 动态
- 职业影响
- 学习资源
- 工具产品

**示例文章 (3条):**
- GPT-5 发布新闻
- 设计师转型指南
- 零基础学习AI路线图

## ✅ 验证安装

```sql
-- 连接数据库
mysql -u root -p

-- 选择数据库
USE hot_ai;

-- 查看所有表
SHOW TABLES;

-- 检查文章分类数量 (应该返回 4)
SELECT COUNT(*) FROM article_categories;

-- 检查文章数量 (应该返回 3)
SELECT COUNT(*) FROM articles;

-- 查看文章分类
SELECT id, name, slug FROM article_categories ORDER BY sort_order;

-- 查看文章列表
SELECT id, title, category, published_at FROM articles;
```

## 🔧 添加新的迁移脚本

创建新文件时使用数字前缀保持顺序:

```bash
# 例如: 添加职业数据
touch migrations/003_seed_profession_data.sql

# 添加学习路径表
touch migrations/004_create_learning_paths.sql
```

## ⚠️ 注意事项

1. **执行顺序**: 必须按数字顺序执行脚本
2. **幂等性**: 所有脚本都使用 `IF NOT EXISTS` 和 `ON DUPLICATE KEY UPDATE`,可安全重复执行
3. **字符集**: 确保使用 utf8mb4 字符集
4. **备份**: 生产环境执行前先备份数据

## 📝 回滚脚本

如需回滚,手动执行:

```sql
-- 删除文章数据
DELETE FROM articles;
DELETE FROM article_categories;

-- 或删除整个数据库(谨慎操作!)
DROP DATABASE hot_ai;
```

## 🐛 故障排查

### 问题: 中文乱码

```bash
mysql -u root -p --default-character-set=utf8mb4 < migrations/002_seed_article_data.sql
```

### 问题: 外键约束错误

确保先执行 `001_create_tables.sql` 创建所有表。

### 问题: 数据库不存在

```sql
CREATE DATABASE hot_ai CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```
