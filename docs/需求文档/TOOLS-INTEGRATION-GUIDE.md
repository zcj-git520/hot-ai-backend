# 工具微服务与前端联调指南

## 概述

本文档描述了如何完成工具微服务（tool-svc）与前端的联调工作。

## 已完成的工作

### 后端 (tool-svc)

1. ✅ **API接口实现**
   - `GET /api/tools/categories` - 获取工具类别列表
   - `GET /api/tools` - 获取工具列表（支持分页和筛选）
   - `GET /api/tools/:slug` - 根据slug获取工具详情
   - `GET /api/tools/id/:id` - 根据ID获取工具详情

2. ✅ **代码修复**
   - 修复了 `tool_repository.go` 中不存在的 `is_free` 字段引用
   - 重新编译并启动了 tool-svc 服务（端口 8004）

3. ✅ **Nginx配置**
   - Nginx已配置代理 `/api/tools` 到 `http://127.0.0.1:8004`

### 前端

1. ✅ **工具列表页面** (`pages/tools/index.vue`)
   - 从API加载真实数据
   - 显示加载状态
   - 显示空数据提示
   - 点击卡片跳转到详情页

2. ✅ **工具详情页面** (`pages/tools/[slug].vue`)
   - 显示工具完整信息
   - 显示定价信息
   - 显示文档链接
   - 显示标签
   - 提供访问官网按钮

3. ✅ **API客户端** (`app/lib/api.ts`)
   - 添加了 `toolApi.getList()` 方法
   - 添加了 `toolApi.getById(slug)` 方法

## 需要执行的步骤

### 1. 插入测试数据

数据库中目前没有工具数据，需要执行SQL脚本插入测试数据。

**方法一：使用提供的SQL文件**

```bash
# 在MySQL命令行中执行
mysql -uroot -pzcj123456 hot_ai < d:\aihot\hot-ai-backend\migrations\test_insert_tools.sql
```

**方法二：手动执行SQL**

打开MySQL客户端，连接到 `hot_ai` 数据库，然后执行以下SQL：

```sql
INSERT INTO `tools` (`id`, `name`, `slug`, `icon`, `description`, `official_url`, `documentation_url`, `pricing`, `pricing_description`, `category_id`, `difficulty`, `rating`, `review_count`, `view_count`, `popularity`, `tags`, `featured`, `status`, `created_by`) VALUES
(1, 'ChatGPT', 'chatgpt', '🤖', 'OpenAI 开发的大型语言模型，能够理解并生成人类语言，适用于写作、编程、学习等多种场景。', 'https://chat.openai.com', 'https://help.openai.com', '{"free":{"available":true,"limit":"每天消息限额","features":["基础对话","基础写作"]},"plus":{"name":"Plus","price":"$20/月","features":["无限对话","GPT-4访问","优先支持"]}}', '免费版每天有消息配额，Plus 版本 $20/月无限额', 1, 'beginner', 4.8, 2300, 15000, 95, '[1,2,3,7]', TRUE, 1, 'system'),
(2, 'Claude', 'claude', '🧠', 'Anthropic 开发的大型语言模型，擅长长文本理解和创作，安全性和准确性优秀。', 'https://claude.ai', 'https://docs.anthropic.com', '{"free":{"available":true,"limit":"有限配额","features":["基础对话","长文本处理"]},"pro":{"name":"Pro","price":"$20/月","features":["无限对话","大文件分析","最高优先级"]}}', '免费版有限配额，Pro 版本 $20/月', 1, 'beginner', 4.8, 1850, 12000, 88, '[1,2,3,7]', TRUE, 1, 'system'),
(3, 'Midjourney', 'midjourney', '🎨', 'AI 图像生成工具，以生成高质量的艺术作品而闻名，支持多种风格和分辨率。', 'https://www.midjourney.com', 'https://docs.midjourney.com', '{"price":"$10/月","features":["无限生成","高分辨率输出","多种风格预设"]}', 'Discord 按生成次数计费，约 $10/月，每月 200 次免费', 2, 'intermediate', 4.8, 1560, 18000, 92, '[4,8,9]', TRUE, 1, 'system'),
(4, 'Stable Diffusion', 'stable-diffusion', '🖼️', '开源的 AI 图像生成模型，可以在本地运行，完全可控，免费使用。', 'https://stability.ai', 'https://platform.stability.ai', '{"free":{"available":true,"limit":"基础功能"},"professional":{"name":"Professional","price":"$20/月","features":["API访问","高优先级","专属资源"]}}', '免费版基础功能，Professional 版本 $20/月', 2, 'intermediate', 4.7, 980, 15000, 85, '[4,7,9]', TRUE, 1, 'system'),
(5, 'GitHub Copilot', 'copilot', '💻', 'GitHub 开发的 AI 编程助手，实时提供代码建议和自动补全。', 'https://github.com/features/copilot', 'https://docs.github.com/en/copilot', '{"price":"$10/月","features":["实时代码建议","多语言支持","安全扫描"]}', '个人版 $10/月，企业版额外费用', 5, 'beginner', 4.7, 1250, 18000, 90, '[6,7,8]', TRUE, 1, 'system'),
(6, 'Cursor', 'cursor', '✏️', '基于 VS Code 的 AI 编程编辑器，集成了多种 AI 功能，代码编写效率高。', 'https://cursor.sh', 'https://docs.cursor.sh', '{"free":{"available":true,"limit":"2000次/月"},"pro":{"name":"Pro","price":"$20/月","features":["无限次数","GPT-4支持","Claude支持"]}}', '免费版每月 2000 次对话，Pro 版本 $20/月', 5, 'beginner', 4.6, 890, 14500, 85, '[6,7,8]', TRUE, 1, 'system');
```

### 2. 验证后端API

执行以下命令验证API是否正常工作：

```powershell
# 测试工具类别接口
curl "http://localhost:8004/api/tools/categories" -UseBasicParsing

# 测试工具列表接口
curl "http://localhost:8004/api/tools?page=1&page_size=5" -UseBasicParsing

# 测试工具详情接口
curl "http://localhost:8004/api/tools/chatgpt" -UseBasicParsing
```

预期结果：应该返回JSON格式的工具数据。

### 3. 启动前端开发服务器

```bash
cd d:\aihot\hot-ai-frontend
pnpm dev
```

前端将在 http://localhost:3000 启动。

### 4. 测试前端页面

1. 打开浏览器访问 http://localhost:3000/tools
2. 应该能看到工具列表页面
3. 点击任意工具卡片，应该能跳转到详情页
4. 详情页应该显示完整的工具信息

## API端点说明

### 通过Nginx代理（推荐）

前端配置使用 `http://localhost/api` 作为基础URL，Nginx会将请求代理到相应的微服务。

- `http://localhost/api/tools/categories` → `http://127.0.0.1:8004/api/tools/categories`
- `http://localhost/api/tools` → `http://127.0.0.1:8004/api/tools`
- `http://localhost/api/tools/:slug` → `http://127.0.0.1:8004/api/tools/:slug`

### 直接访问微服务

也可以直接访问微服务进行测试：

- `http://localhost:8004/api/tools/categories`
- `http://localhost:8004/api/tools`
- `http://localhost:8004/api/tools/:slug`

## 常见问题

### 1. 前端显示"暂无工具数据"

**原因**: 数据库中没有工具数据

**解决**: 执行步骤1中的SQL脚本插入测试数据

### 2. API返回404错误

**原因**: Nginx未运行或配置不正确

**解决**: 
- 检查Nginx是否正在运行
- 确认 `nginx.conf` 中包含 `/api/tools` 的配置
- 重启Nginx

### 3. 前端无法连接到后端

**原因**: 后端服务未启动

**解决**:
```powershell
# 检查tool-svc是否在运行
netstat -ano | findstr ":8004"

# 如果没有运行，重新启动
cd d:\aihot\hot-ai-backend
.\bin\tool-svc.exe
```

### 4. CORS错误

**原因**: 跨域请求被阻止

**解决**: 
- 确保Nginx配置中包含CORS头
- 或者直接在开发时使用微服务的地址

## 下一步优化建议

1. **添加搜索功能**: 在工具列表页面添加搜索框
2. **添加筛选功能**: 按类别、难度等筛选工具
3. **添加工具评测**: 允许用户对工具进行评分和评论
4. **添加收藏功能**: 用户可以收藏喜欢的工具
5. **优化工具详情**: 添加更多相关信息，如使用教程、示例等

## 技术栈

- **后端**: Go + go-zero框架 + MySQL
- **前端**: Nuxt 3 + Vue 3 + TypeScript
- **代理**: Nginx
- **数据库**: MySQL 8.0+

## 联系与支持

如有问题，请查看：
- 后端日志: `d:\aihot\hot-ai-backend\logs\tool-svc.log`
- API文档: `d:\aihot\hot-ai-backend\docs\api\tools-api.md`
