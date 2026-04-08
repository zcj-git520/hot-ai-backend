-- ======================================================
-- AI 热点追踪平台 - 学习路径初始数据
-- 版本：v1.0
-- 日期：2026-04-08
-- ======================================================

USE `hot_ai`;

-- ============================================
-- 1. 插入学习路径数据（6 条）
-- ============================================

INSERT INTO `learning_paths` (`id`, `title`, `slug`, `icon`, `description`, `difficulty`, `level_label`, `learning_goals`, `target_audience`, `estimated_days`, `estimated_hours`, `chapter_count`, `student_count`, `cover_image`, `is_featured`, `is_active`, `sort_order`, `status`, `published_at`) VALUES
(1, '零基础入门 AI', 'zero-to-ai', '🌱', '从零开始，系统学习 AI 基础知识，掌握常用工具和实践技能', 'beginner', '入门', '[\"理解 AI 基本概念和工作原理\", \"掌握主流 AI 工具的使用方法\", \"能够使用 AI 辅助日常工作\", \"了解 AI 在各行业的应用场景\"]', '[\"AI 完全零基础的用户\", \"希望系统入门 AI 的职场人士\", \"对 AI 感兴趣的学生\", \"想要提升工作效率的办公人员\"]', 30, 60, 12, 12500, NULL, 1, 1, 1, 1, NOW()),
(2, '提示词工程进阶', 'prompt-engineering', '✍️', '掌握高级提示词技巧，提升与大模型的交互效率和质量', 'intermediate', '进阶', '[\"掌握提示词工程的核心原理\", \"学会编写高效的结构化提示词\", \"理解不同场景下的提示词策略\", \"能够优化和调试提示词\"]', '[\"有一定 AI 使用经验的用户\", \"文案策划和内容创作者\", \"需要频繁使用 AI 工具的职场人士\"]', 21, 42, 8, 8300, NULL, 1, 1, 2, 1, NOW()),
(3, 'AI 编程助手实战', 'ai-programming', '💻', '学习使用 Cursor、GitHub Copilot 等工具提升开发效率', 'intermediate', '进阶', '[\"掌握 AI 编程助手的核心功能\", \"学会使用 AI 生成和优化代码\", \"理解 AI 辅助调试和测试的方法\", \"建立高效的 AI 辅助开发工作流\"]', '[\"程序员和软件开发人员\", \"数据分析师和工程师\", \"想要提升编码效率的开发者\"]', 28, 56, 10, 15200, NULL, 1, 1, 3, 1, NOW()),
(4, 'AI 绘画与视觉设计', 'ai-painting', '🎨', '掌握 Midjourney、Stable Diffusion 等 AI 绘画工具', 'beginner', '入门', '[\"了解 AI 绘画工具的基本原理\", \"掌握主流 AI 绘画工具的使用\", \"学会提示词生成理想图像\", \"能够完成商业级 AI 绘画作品\"]', '[\"设计师和视觉创作者\", \"对 AI 绘画感兴趣的用户\", \"想要尝试数字艺术创作的初学者\"]', 35, 70, 14, 9700, NULL, 1, 1, 4, 1, NOW()),
(5, 'AI 应用开发', 'ai-app-development', '🚀', '学习构建基于大模型的实际应用和产品', 'advanced', '高级', '[\"掌握 LLM API 调用和集成\", \"学会构建 AI 驱动的应用程序\", \"理解 RAG 和 Agent 架构\", \"能够部署和运维 AI 应用\"]', '[\"有一定编程基础的开发者\", \"产品经理和技术负责人\", \"想要进入 AI 领域的技术人员\"]', 45, 90, 16, 5600, NULL, 1, 1, 5, 1, NOW()),
(6, '数据分析与 AI', 'data-analysis-ai', '📊', '运用 AI 技术提升数据分析能力和商业洞察', 'intermediate', '进阶', '[\"掌握 AI 辅助数据分析的方法\", \"学会使用 AI 进行数据可视化\", \"理解 AI 在商业智能中的应用\", \"能够构建 AI 驱动的数据分析报告\"]', '[\"数据分析师和商业分析师\", \"市场营销人员\", \"需要数据驱动的决策者\"]', 40, 80, 13, 7100, NULL, 0, 1, 6, 1, NOW());

-- ============================================
-- 2. 插入章节数据
-- ============================================

-- 路径1：零基础入门 AI 的章节 (path_id = 1)
INSERT INTO `path_chapters` (`path_id`, `title`, `slug`, `description`, `content_type`, `content`, `video_url`, `external_links`, `estimated_hours`, `order_index`, `is_free`, `status`) VALUES
(1, 'AI 时代背景与趋势', 'ai-background-trends', '了解 AI 发展历程和当前热点', 'article', '# AI 时代背景与趋势\n\n## 什么是人工智能\n\n人工智能（Artificial Intelligence，简称 AI）是计算机科学的一个分支...', NULL, NULL, 2, 1, 1, 1),
(1, '大语言模型基础概念', 'llm-basics', '理解 GPT、Claude 等模型原理', 'article', '# 大语言模型基础概念\n\n## 什么是大语言模型\n\n大语言模型（Large Language Model，LLM）是一种基于深度学习...', NULL, NULL, 3, 2, 1, 1),
(1, '提示词工程入门', 'prompt-basics', '学习基础提示词编写技巧', 'article', '# 提示词工程入门\n\n## 什么是提示词\n\n提示词（Prompt）是用户输入给 AI 模型的指令或问题...', NULL, NULL, 4, 3, 1, 1),
(1, 'ChatGPT 实战应用', 'chatgpt-practice', '掌握 ChatGPT 的常用功能', 'practice', '# ChatGPT 实战应用\n\n## 实践任务\n\n1. 使用 ChatGPT 完成一篇博客文章\n2. 让 ChatGPT 帮你制定学习计划...', NULL, NULL, 3, 4, 1, 1),
(1, 'AI 图像生成工具入门', 'ai-image-basics', '学习 Midjourney 基础使用', 'article', '# AI 图像生成工具入门\n\n## Midjourney 基础\n\nMidjourney 是目前最流行的 AI 绘画工具之一...', NULL, NULL, 4, 5, 1, 1),
(1, 'AI 语音工具使用', 'ai-voice-tools', '了解语音合成和识别工具', 'article', '# AI 语音工具使用\n\n## 语音合成技术\n\n语音合成（Text-to-Speech，TTS）将文本转换为自然语音...', NULL, NULL, 2, 6, 1, 1),
(1, 'AI 写作助手实战', 'ai-writing-practice', '使用 AI 提升写作效率', 'practice', '# AI 写作助手实战\n\n## 实践任务\n\n1. 使用 AI 生成文章大纲\n2. 让 AI 优化文章表达...', NULL, NULL, 3, 7, 1, 1),
(1, 'AI 办公提效技巧', 'ai-productivity', 'AI 在日常办公中的应用', 'article', '# AI 办公提效技巧\n\n## AI 辅助办公\n\nAI 可以帮助我们处理各种办公任务...', NULL, NULL, 2, 8, 1, 1),
(1, 'AI 工具选型指南', 'ai-tool-selection', '如何选择适合的 AI 工具', 'article', '# AI 工具选型指南\n\n## 工具选择原则\n\n1. 根据需求选择\n2. 考虑成本因素...', NULL, NULL, 2, 9, 1, 1),
(1, 'AI 伦理与安全', 'ai-ethics-safety', '了解 AI 使用的注意事项', 'article', '# AI 伦理与安全\n\n## AI 伦理问题\n\n随着 AI 技术的普及，伦理问题日益重要...', NULL, NULL, 2, 10, 1, 1),
(1, '学习资源汇总', 'ai-learning-resources', '推荐进一步学习的资源', 'external', '# 学习资源汇总\n\n## 推荐书籍\n\n1. 《人工智能：一种现代方法》\n2. 《深度学习》...', NULL, '[\"https://www.coursera.org/learn/ai-for-everyone\", \"https://www.deeplearning.ai/\"]', 1, 11, 1, 1),
(1, '综合实践项目', 'ai-capstone-project', '完成一个完整的 AI 应用项目', 'practice', '# 综合实践项目\n\n## 项目要求\n\n完成一个完整的 AI 应用项目...', NULL, NULL, 4, 12, 1, 1);

-- 路径2：提示词工程进阶的章节 (path_id = 2)
INSERT INTO `path_chapters` (`path_id`, `title`, `slug`, `description`, `content_type`, `content`, `video_url`, `external_links`, `estimated_hours`, `order_index`, `is_free`, `status`) VALUES
(2, '提示词工程原理', 'prompt-engineering-principles', '理解提示词如何影响模型输出', 'article', '# 提示词工程原理\n\n## 提示词的作用\n\n提示词直接影响模型的输出质量和准确性...', NULL, NULL, 3, 1, 1, 1),
(2, '结构化提示词设计', 'structured-prompts', '学习编写结构化的提示词', 'article', '# 结构化提示词设计\n\n## 设计原则\n\n1. 清晰明确\n2. 提供上下文\n3. 设定约束条件...', NULL, NULL, 4, 2, 1, 1),
(2, 'Few-Shot 学习技巧', 'few-shot-learning', '掌握少样本学习的方法', 'article', '# Few-Shot 学习技巧\n\n## 什么是 Few-Shot\n\nFew-Shot Learning 通过提供少量示例来引导模型...', NULL, NULL, 3, 3, 1, 1),
(2, 'Chain of Thought 思维链', 'chain-of-thought', '学习思维链提示技巧', 'article', '# Chain of Thought 思维链\n\n## 思维链原理\n\n思维链通过引导模型展示推理过程来提高准确性...', NULL, NULL, 4, 4, 1, 1),
(2, '提示词优化与调试', 'prompt-optimization', '如何优化和调试提示词', 'practice', '# 提示词优化与调试\n\n## 实践任务\n\n1. 优化一个低质量的提示词\n2. A/B 测试不同的提示词版本...', NULL, NULL, 4, 5, 1, 1),
(2, '场景化提示词实战', 'scenario-prompts', '不同场景下的提示词策略', 'practice', '# 场景化提示词实战\n\n## 场景练习\n\n1. 客服场景\n2. 教育场景\n3. 创意场景...', NULL, NULL, 5, 6, 1, 1),
(2, '提示词模板库构建', 'prompt-template-library', '建立自己的提示词库', 'practice', '# 提示词模板库构建\n\n## 实践任务\n\n建立自己的提示词模板库，分类整理常用提示词...', NULL, NULL, 3, 7, 1, 1),
(2, '高级技巧综合实践', 'advanced-prompts-project', '综合运用高级提示词技巧', 'practice', '# 高级技巧综合实践\n\n## 项目任务\n\n综合运用所有技巧，完成一个复杂提示词项目...', NULL, NULL, 6, 8, 1, 1);

-- 路径3：AI 编程助手实战的章节 (path_id = 3)
INSERT INTO `path_chapters` (`path_id`, `title`, `slug`, `description`, `content_type`, `content`, `video_url`, `external_links`, `estimated_hours`, `order_index`, `is_free`, `status`) VALUES
(3, 'AI 编程助手概览', 'ai-coding-overview', '了解主流 AI 编程工具', 'article', '# AI 编程助手概览\n\n## 主流工具\n\n1. GitHub Copilot\n2. Cursor\n3. CodeWhisperer...', NULL, NULL, 2, 1, 1, 1),
(3, 'Cursor 基础使用', 'cursor-basics', '掌握 Cursor 的核心功能', 'article', '# Cursor 基础使用\n\n## 功能介绍\n\nCursor 是一款 AI 驱动的代码编辑器...', NULL, NULL, 3, 2, 1, 1),
(3, 'GitHub Copilot 实战', 'copilot-practice', '学习使用 Copilot 编写代码', 'practice', '# GitHub Copilot 实战\n\n## 实践任务\n\n1. 使用 Copilot 实现一个函数\n2. 让 Copilot 生成测试用例...', NULL, NULL, 4, 3, 1, 1),
(3, 'AI 辅助代码审查', 'ai-code-review', '使用 AI 进行代码审查', 'article', '# AI 辅助代码审查\n\n## 审查要点\n\n1. 代码规范\n2. 潜在 bug\n3. 性能优化建议...', NULL, NULL, 3, 4, 1, 1),
(3, 'AI 生成单元测试', 'ai-unit-tests', '让 AI 帮你写测试代码', 'practice', '# AI 生成单元测试\n\n## 实践任务\n\n使用 AI 为现有代码生成单元测试...', NULL, NULL, 4, 5, 1, 1),
(3, 'AI 辅助调试技巧', 'ai-debugging', '使用 AI 快速定位和修复 bug', 'article', '# AI 辅助调试技巧\n\n## 调试方法\n\n1. 让 AI 分析错误日志\n2. 使用 AI 解释异常行为...', NULL, NULL, 3, 6, 1, 1),
(3, 'AI 重构代码', 'ai-refactoring', '让 AI 帮你优化代码结构', 'practice', '# AI 重构代码\n\n## 实践任务\n\n重构一个复杂函数，提高代码可读性...', NULL, NULL, 4, 7, 1, 1),
(3, 'AI 文档生成', 'ai-documentation', '使用 AI 生成代码文档', 'article', '# AI 文档生成\n\n## 文档类型\n\n1. API 文档\n2. README\n3. 注释生成...', NULL, NULL, 2, 8, 1, 1),
(3, 'AI 工作流搭建', 'ai-workflow', '建立高效的 AI 辅助开发流程', 'practice', '# AI 工作流搭建\n\n## 实践任务\n\n建立一套完整的 AI 辅助开发工作流...', NULL, NULL, 4, 9, 1, 1),
(3, '综合实战项目', 'ai-coding-project', '使用 AI 完成完整项目', 'practice', '# 综合实战项目\n\n## 项目要求\n\n使用 AI 编程助手完成一个完整项目...', NULL, NULL, 7, 10, 1, 1);

-- 路径4：AI 绘画与视觉设计的章节 (path_id = 4)
INSERT INTO `path_chapters` (`path_id`, `title`, `slug`, `description`, `content_type`, `content`, `video_url`, `external_links`, `estimated_hours`, `order_index`, `is_free`, `status`) VALUES
(4, 'AI 绘画工具概览', 'ai-painting-overview', '了解主流 AI 绘画工具', 'article', '# AI 绘画工具概览\n\n## 主流工具\n\n1. Midjourney\n2. Stable Diffusion\n3. DALL-E 3...', NULL, NULL, 2, 1, 1, 1),
(4, 'Midjourney 基础入门', 'midjourney-basics', '学习 Midjourney 基本命令', 'article', '# Midjourney 基础入门\n\n## 基本命令\n\n1. /imagine - 生成图像\n2. /blend - 混合图像...', NULL, NULL, 4, 2, 1, 1),
(4, 'Midjourney 参数详解', 'midjourney-parameters', '掌握 MJ 的各种参数', 'article', '# Midjourney 参数详解\n\n## 常用参数\n\n1. --ar 设置宽高比\n2. --v 版本选择\n3. --s 风格化...', NULL, NULL, 4, 3, 1, 1),
(4, 'Stable Diffusion 入门', 'stable-diffusion-basics', '了解 SD 基本原理和使用', 'article', '# Stable Diffusion 入门\n\n## 基本原理\n\nStable Diffusion 是一种潜在扩散模型...', NULL, NULL, 4, 4, 1, 1),
(4, '提示词与风格控制', 'prompt-style-control', '学习控制图像风格和细节', 'article', '# 提示词与风格控制\n\n## 风格控制技巧\n\n1. 使用艺术家风格\n2. 指定媒介材质\n3. 控制光线氛围...', NULL, NULL, 4, 5, 1, 1),
(4, 'ControlNet 使用技巧', 'controlnet-tips', '掌握 ControlNet 精准控制', 'practice', '# ControlNet 使用技巧\n\n## 实践任务\n\n使用 ControlNet 精准控制图像生成...', NULL, NULL, 5, 6, 1, 1),
(4, 'LoRA 模型训练', 'lora-training', '训练自己的风格模型', 'practice', '# LoRA 模型训练\n\n## 实践任务\n\n训练一个 LoRA 模型，学习特定风格...', NULL, NULL, 6, 7, 1, 1),
(4, '图像后期处理', 'image-post-processing', 'AI 生成图像的后期优化', 'article', '# 图像后期处理\n\n## 优化方法\n\n1. 分辨率提升\n2. 细节修复\n3. 色彩调整...', NULL, NULL, 3, 8, 1, 1),
(4, '商业应用场景', 'commercial-applications', 'AI 绘画的商业应用', 'article', '# 商业应用场景\n\n## 应用领域\n\n1. 游戏美术\n2. 广告设计\n3. 电商展示...', NULL, NULL, 3, 9, 1, 1),
(4, '版权与伦理问题', 'copyright-ethics', '了解 AI 绘画的版权问题', 'article', '# 版权与伦理问题\n\n## 注意事项\n\n1. 版权归属\n2. 商业使用限制\n3. 伦理考量...', NULL, NULL, 2, 10, 1, 1),
(4, '作品集打造', 'portfolio-building', '打造专业的 AI 艺术作品集', 'practice', '# 作品集打造\n\n## 实践任务\n\n创建个人 AI 艺术作品集...', NULL, NULL, 4, 11, 1, 1),
(4, '综合创作项目', 'creative-project', '完成一个完整的创作项目', 'practice', '# 综合创作项目\n\n## 项目要求\n\n完成一个完整的 AI 艺术创作项目...', NULL, NULL, 5, 12, 1, 1),
(4, '进阶技巧探索', 'advanced-techniques', '探索更多高级技巧', 'external', '# 进阶技巧探索\n\n## 学习资源\n\n1. 社区教程\n2. 论文解读\n3. 前沿技术...', NULL, '[\"https://stable-diffusion-art.com/\", \"https://www.midjourney.com/showcase\"]', 3, 13, 1, 1),
(4, '社区资源分享', 'community-resources', '推荐学习和交流社区', 'external', '# 社区资源分享\n\n## 推荐社区\n\n1. Reddit r/StableDiffusion\n2. Discord 社区\n3. B站教程...', NULL, '[\"https://discord.gg/midjourney\", \"https://www.reddit.com/r/StableDiffusion/\"]', 2, 14, 1, 1);

-- 路径5：AI 应用开发的章节 (path_id = 5)
INSERT INTO `path_chapters` (`path_id`, `title`, `slug`, `description`, `content_type`, `content`, `video_url`, `external_links`, `estimated_hours`, `order_index`, `is_free`, `status`) VALUES
(5, 'LLM API 基础', 'llm-api-basics', '学习调用大模型 API', 'article', '# LLM API 基础\n\n## API 调用示例\n\n```python\nimport openai\n\nresponse = openai.ChatCompletion.create(\n    model=\"gpt-4\",\n    messages=[{\"role\": \"user\", \"content\": \"Hello\"}]\n)\n```', NULL, NULL, 3, 1, 1, 1),
(5, 'LangChain 框架入门', 'langchain-intro', '掌握 LangChain 核心概念', 'article', '# LangChain 框架入门\n\n## 核心概念\n\n1. Chains\n2. Agents\n3. Tools\n4. Memory...', NULL, NULL, 4, 2, 1, 1),
(5, 'Prompt 模板管理', 'prompt-template-mgmt', '管理和优化提示词模板', 'article', '# Prompt 模板管理\n\n## 模板设计\n\n创建可复用的提示词模板...', NULL, NULL, 3, 3, 1, 1),
(5, 'Chain 链式调用', 'chain-calls', '构建复杂的调用链', 'practice', '# Chain 链式调用\n\n## 实践任务\n\n构建一个多步骤的 Chain 处理流程...', NULL, NULL, 4, 4, 1, 1),
(5, 'Agent 智能体开发', 'agent-development', '创建自主决策的 AI 智能体', 'practice', '# Agent 智能体开发\n\n## 实践任务\n\n创建一个能够自主决策和行动的 Agent...', NULL, NULL, 6, 5, 1, 1),
(5, 'RAG 检索增强生成', 'rag-implementation', '实现 RAG 架构', 'practice', '# RAG 检索增强生成\n\n## 实践任务\n\n实现一个完整的 RAG 系统...', NULL, NULL, 6, 6, 1, 1),
(5, '向量数据库使用', 'vector-database', '学习向量数据库的用法', 'article', '# 向量数据库使用\n\n## 主流数据库\n\n1. Pinecone\n2. Weaviate\n3. Qdrant\n4. Milvus...', NULL, NULL, 4, 7, 1, 1),
(5, '知识库构建', 'knowledge-base', '构建和管理知识库', 'practice', '# 知识库构建\n\n## 实践任务\n\n构建一个领域知识库...', NULL, NULL, 5, 8, 1, 1),
(5, 'AI 应用前端集成', 'frontend-integration', '将 AI 功能集成到前端', 'practice', '# AI 应用前端集成\n\n## 实践任务\n\n在前端应用中集成 AI 功能...', NULL, NULL, 4, 9, 1, 1),
(5, '流式响应处理', 'streaming-responses', '处理流式 API 响应', 'article', '# 流式响应处理\n\n## 实现方法\n\n处理 SSE (Server-Sent Events) 流式响应...', NULL, NULL, 3, 10, 1, 1),
(5, '错误处理与重试', 'error-handling', '健壮的错误处理机制', 'article', '# 错误处理与重试\n\n## 最佳实践\n\n1. 指数退避重试\n2. 错误分类处理\n3. 降级策略...', NULL, NULL, 3, 11, 1, 1),
(5, '性能优化技巧', 'performance-optimization', '优化 AI 应用性能', 'article', '# 性能优化技巧\n\n## 优化方法\n\n1. 缓存策略\n2. 并发控制\n3. 批处理...', NULL, NULL, 4, 12, 1, 1),
(5, '部署与运维', 'deployment-ops', 'AI 应用的部署和监控', 'article', '# 部署与运维\n\n## 部署方案\n\n1. Docker 容器化\n2. Kubernetes 编排\n3. 监控告警...', NULL, NULL, 4, 13, 1, 1),
(5, '成本优化策略', 'cost-optimization', '降低 API 调用成本', 'article', '# 成本优化策略\n\n## 优化技巧\n\n1. 模型选择\n2. 缓存机制\n3. 请求合并...', NULL, NULL, 3, 14, 1, 1),
(5, '安全与合规', 'security-compliance', 'AI 应用的安全考虑', 'article', '# 安全与合规\n\n## 安全要点\n\n1. API 密钥管理\n2. 输入验证\n3. 数据隐私...', NULL, NULL, 3, 15, 1, 1),
(5, '综合项目实战', 'capstone-project', '完成一个完整的 AI 应用', 'practice', '# 综合项目实战\n\n## 项目要求\n\n完成一个端到端的 AI 应用项目...', NULL, NULL, 12, 16, 1, 1);

-- 路径6：数据分析与 AI 的章节 (path_id = 6)
INSERT INTO `path_chapters` (`path_id`, `title`, `slug`, `description`, `content_type`, `content`, `video_url`, `external_links`, `estimated_hours`, `order_index`, `is_free`, `status`) VALUES
(6, 'AI 数据分析概览', 'ai-data-analysis-overview', '了解 AI 在数据分析中的应用', 'article', '# AI 数据分析概览\n\n## 应用场景\n\n1. 数据清洗\n2. 探索性分析\n3. 预测建模...', NULL, NULL, 2, 1, 1, 1),
(6, '数据预处理自动化', 'data-preprocessing', '使用 AI 自动化数据清洗', 'article', '# 数据预处理自动化\n\n## 清洗方法\n\n1. 缺失值处理\n2. 异常值检测\n3. 数据标准化...', NULL, NULL, 3, 2, 1, 1),
(6, 'AI 辅助数据探索', 'ai-data-exploration', '让 AI 帮你探索数据', 'practice', '# AI 辅助数据探索\n\n## 实践任务\n\n使用 AI 辅助探索数据集...', NULL, NULL, 4, 3, 1, 1),
(6, '智能数据可视化', 'ai-visualization', 'AI 生成可视化图表', 'practice', '# 智能数据可视化\n\n## 实践任务\n\n让 AI 自动生成合适的可视化图表...', NULL, NULL, 4, 4, 1, 1),
(6, '自然语言查询数据', 'nl-data-query', '用自然语言查询数据', 'article', '# 自然语言查询数据\n\n## 实现技术\n\n1. Text-to-SQL\n2. NL-to-Query...', NULL, NULL, 3, 5, 1, 1),
(6, 'AI 异常检测', 'ai-anomaly-detection', '使用 AI 检测数据异常', 'practice', '# AI 异常检测\n\n## 实践任务\n\n实现异常检测模型...', NULL, NULL, 4, 6, 1, 1),
(6, '预测分析入门', 'predictive-analytics', 'AI 驱动的预测分析', 'article', '# 预测分析入门\n\n## 预测方法\n\n1. 时间序列预测\n2. 回归分析\n3. 分类预测...', NULL, NULL, 4, 7, 1, 1),
(6, '文本数据分析', 'text-data-analysis', '使用 NLP 分析文本数据', 'practice', '# 文本数据分析\n\n## 实践任务\n\n分析文本数据的主题和情感...', NULL, NULL, 5, 8, 1, 1),
(6, 'AI 报告生成', 'ai-report-generation', '自动生成数据分析报告', 'practice', '# AI 报告生成\n\n## 实践任务\n\n自动生成数据分析报告...', NULL, NULL, 4, 9, 1, 1),
(6, '商业智能应用', 'bi-applications', 'AI 在 BI 中的应用', 'article', '# 商业智能应用\n\n## 应用场景\n\n1. 销售分析\n2. 用户行为分析\n3. 运营监控...', NULL, NULL, 3, 10, 1, 1),
(6, 'Dashboard 自动化', 'dashboard-automation', '自动化仪表板构建', 'practice', '# Dashboard 自动化\n\n## 实践任务\n\n自动构建数据仪表板...', NULL, NULL, 4, 11, 1, 1),
(6, '数据故事讲述', 'data-storytelling', '用 AI 讲述数据故事', 'article', '# 数据故事讲述\n\n## 技巧要点\n\n1. 发现洞察\n2. 构建叙事\n3. 可视化呈现...', NULL, NULL, 3, 12, 1, 1),
(6, '综合实战项目', 'analysis-project', '完成一个完整的数据分析项目', 'practice', '# 综合实战项目\n\n## 项目要求\n\n完成一个端到端的数据分析项目...', NULL, NULL, 7, 13, 1, 1);

-- ============================================
-- 3. 插入学习路径管理数据
-- ============================================

INSERT INTO `learning_path_management` (`path_id`, `view_count`, `start_count`, `complete_count`, `favorite_count`, `meta_title`, `meta_description`, `meta_keywords`, `reviewer_id`, `reviewed_at`, `review_notes`) VALUES
(1, 125000, 12500, 3200, 8500, '零基础入门 AI - 30 天系统学习 AI 基础知识 | AI 热点追踪', '从零开始学习 AI，掌握大语言模型、提示词工程、AI 工具使用等核心技能，适合完全零基础的用户。', 'AI 入门,零基础学 AI,大语言模型,提示词工程,AI 工具', NULL, NOW(), NULL),
(2, 83000, 8300, 2100, 5200, '提示词工程进阶 - 掌握高级提示词技巧 | AI 热点追踪', '深入学习提示词工程，掌握结构化提示、思维链、Few-Shot 等高级技巧，提升与大模型交互效率。', '提示词工程,Prompt Engineering,高级提示词,思维链', NULL, NOW(), NULL),
(3, 152000, 15200, 4800, 11000, 'AI 编程助手实战 - 使用 Cursor 和 Copilot 提升开发效率 | AI 热点追踪', '学习使用 AI 编程助手，包括 Cursor、GitHub Copilot 等工具，大幅提升代码开发效率。', 'AI 编程,Cursor,GitHub Copilot,AI 辅助开发', NULL, NOW(), NULL),
(4, 97000, 9700, 2500, 6800, 'AI 绘画与视觉设计 - 掌握 Midjourney 和 Stable Diffusion | AI 热点追踪', '系统学习 AI 绘画工具，从 Midjourney 入门到 Stable Diffusion 进阶，打造专业 AI 艺术作品。', 'AI 绘画,Midjourney,Stable Diffusion,AI 艺术', NULL, NOW(), NULL),
(5, 56000, 5600, 1200, 3400, 'AI 应用开发 - 构建基于大模型的实际应用 | AI 热点追踪', '学习使用 LangChain、RAG、Agent 等技术，构建完整的 AI 驱动应用程序。', 'AI 应用开发,LangChain,RAG,Agent,大模型应用', NULL, NOW(), NULL),
(6, 71000, 7100, 1800, 4500, '数据分析与 AI - 运用 AI 提升数据分析能力 | AI 热点追踪', '掌握 AI 辅助数据分析的方法，实现数据预处理、可视化、报告生成的自动化。', 'AI 数据分析,智能可视化,AI 报告生成,商业智能', NULL, NOW(), NULL);

-- ============================================
-- 4. 验证数据
-- ============================================

-- 查看学习路径总数
SELECT COUNT(*) as '总路径数' FROM `learning_paths`;

-- 查看各难度等级分布
SELECT `difficulty`, `level_label`, COUNT(*) as '数量' FROM `learning_paths` GROUP BY `difficulty`, `level_label`;

-- 查看章节总数
SELECT COUNT(*) as '总章节数' FROM `path_chapters`;

-- 查看各路径章节数
SELECT 
    p.`id`,
    p.`title`, 
    COUNT(c.`id`) as '章节数',
    p.`chapter_count` as '预期章节数'
FROM `learning_paths` p
LEFT JOIN `path_chapters` c ON p.`id` = c.`path_id`
GROUP BY p.`id`, p.`title`, p.`chapter_count`;

-- 查看管理数据
SELECT 
    p.`title`,
    m.`view_count`,
    m.`start_count`,
    m.`complete_count`,
    m.`favorite_count`
FROM `learning_path_management` m
JOIN `learning_paths` p ON m.`path_id` = p.`id`;

-- ============================================
-- 5. 数据清理脚本（可选）
-- ============================================

-- 如果需要清空所有数据重新导入，使用以下命令
-- SET FOREIGN_KEY_CHECKS = 0;
-- DELETE FROM `learning_path_management`;
-- DELETE FROM `path_chapters`;
-- DELETE FROM `learning_paths`;
-- ALTER TABLE `learning_paths` AUTO_INCREMENT = 1;
-- ALTER TABLE `path_chapters` AUTO_INCREMENT = 1;
-- ALTER TABLE `learning_path_management` AUTO_INCREMENT = 1;
-- SET FOREIGN_KEY_CHECKS = 1;