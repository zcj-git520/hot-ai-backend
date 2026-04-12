-- ============================================================================
-- 工具库模块 - 初始化样例数据
-- ============================================================================
-- 版本: v1.0
-- 创建日期: 2026-04-11
-- 说明: 工具库模块的初始数据，用于测试和演示
-- ============================================================================

-- ============================================================================
-- 1. 工具类别数据 (7大类)
-- ============================================================================

INSERT INTO `tool_categories` (`id`, `name`, `slug`, `icon`, `description`, `sort_order`, `featured`, `status`) VALUES
(1, '写作类', 'writing', '✍️', '用于写作、文案创作的 AI 工具', 1, TRUE, 1),
(2, '图像类', 'image', '🎨', '用于图像生成、编辑的 AI 工具', 2, TRUE, 1),
(3, '视频类', 'video', '🎬', '用于视频生成、编辑的 AI 工具', 3, TRUE, 1),
(4, '音频类', 'audio', '🔊', '用于音频生成、编辑的 AI 工具', 4, TRUE, 1),
(5, '代码类', 'code', '💻', '用于编程、代码相关的 AI 工具', 5, TRUE, 1),
(6, '办公类', 'office', '📊', '用于办公、文档处理的 AI 工具', 6, TRUE, 1),
(7, '其他类', 'other', '🔧', '其他 AI 工具', 7, FALSE, 1);

-- ============================================================================
-- 2. 提示词模板分类数据 (6个类别)
-- ============================================================================

INSERT INTO `prompt_template_categories` (`id`, `name`, `slug`, `description`, `icon`, `sort_order`, `featured`, `status`) VALUES
(1, '写作类', 'writing', '用于写作的提示词模板', '✍️', 1, TRUE, 1),
(2, '代码类', 'code', '用于编程的提示词模板', '💻', 2, TRUE, 1),
(3, '设计类', 'design', '用于设计的提示词模板', '🎨', 3, TRUE, 1),
(4, '营销类', 'marketing', '用于营销的提示词模板', '📢', 4, TRUE, 1),
(5, '学习类', 'learning', '用于学习的提示词模板', '📚', 5, FALSE, 1),
(6, '翻译类', 'translation', '用于翻译的提示词模板', '🌐', 6, FALSE, 1);

-- ============================================================================
-- 3. 工具标签数据
-- ============================================================================

INSERT INTO `tool_tags` (`id`, `name`, `slug`, `color`) VALUES
(1, '大语言模型', 'llm', '#3B82F6'),
(2, '对话', 'chat', '#10B981'),
(3, '写作', 'writing', '#8B5CF6'),
(4, '图像生成', 'image-gen', '#EC4899'),
(5, '视频生成', 'video-gen', '#F59E0B'),
(6, '编程辅助', 'coding', '#06B6D4'),
(7, '免费', 'free', '#22C55E'),
(8, '付费', 'paid', '#EF4444'),
(9, '开源', 'open-source', '#6B7280'),
(10, '国内', 'chinese', '#3B82F6');

-- ============================================================================
-- 4. 工具数据 - 主流AI工具 (20个)
-- ============================================================================

-- 写作类工具 (6个)
INSERT INTO `tools` (`id`, `name`, `slug`, `icon`, `description`, `official_url`, `documentation_url`, `pricing`, `pricing_description`, `category_id`, `difficulty`, `rating`, `review_count`, `view_count`, `popularity`, `tags`, `featured`, `status`, `created_by`) VALUES
(1, 'ChatGPT', 'chatgpt', '🤖', 'OpenAI 开发的大型语言模型，能够理解并生成人类语言，适用于写作、编程、学习等多种场景。', 'https://chat.openai.com', 'https://help.openai.com', '{"free":{"available":true,"limit":"每天消息限额","features":["基础对话","基础写作"]},"plus":{"name":"Plus","price":"$20/月","features":["无限对话","GPT-4访问","优先支持"]}}', '免费版每天有消息配额，Plus 版本 $20/月无限额', 1, 'beginner', 4.8, 2300, 15000, 95, '[1,2,3,7]', TRUE, 1, 'system'),
(2, 'Claude', 'claude', '🧠', 'Anthropic 开发的大型语言模型，擅长长文本理解和创作，安全性和准确性优秀。', 'https://claude.ai', 'https://docs.anthropic.com', '{"free":{"available":true,"limit":"有限配额","features":["基础对话","长文本处理"]},"pro":{"name":"Pro","price":"$20/月","features":["无限对话","大文件分析","最高优先级"]}}', '免费版有限配额，Pro 版本 $20/月', 1, 'beginner', 4.8, 1850, 12000, 88, '[1,2,3,7]', TRUE, 1, 'system'),
(3, 'Notion AI', 'notion-ai', '📝', 'Notion 集成的 AI 助手，帮助用户快速写作、总结、翻译和编辑文档。', 'https://www.notion.so', 'https://www.notion.so/product/ai', '{"price":"$10/月","features":["Notion集成的AI写作助手","文档自动生成和润色","智能搜索和总结"]}', 'Notion 套餐内包含，单独 $10/月', 1, 'beginner', 4.7, 980, 8500, 75, '[1,2,3,6,7]', TRUE, 1, 'system'),
(4, 'Jasper', 'jasper', '✍️', '专为营销文案设计的 AI 写作工具，提供超过 50 种模板。', 'https://jasper.ai', 'https://www.jasper.ai', '{"free":{"available":false,"trial":"5天试用"},"creator":{"name":"Creator","price":"$49/月","features":["50+模板","多语言支持"]}}', '免费试用 5 天，Creator 版本 $49/月', 1, 'intermediate', 4.5, 650, 7200, 70, '[3,8]', FALSE, 1, 'system'),
(5, 'Copy.ai', 'copy-ai', '📋', '快速生成营销文案、社交媒体内容、博客文章的 AI 写作工具。', 'https://copy.ai', 'https://www.copy.ai', '{"free":{"available":true,"limit":"每周10次"},"team":{"name":"Team","price":"$36/月/用户","features":["无限次数","团队协作"]}}', '免费版每周 10 次，Team 版本 $36/月/用户', 1, 'beginner', 4.3, 420, 6800, 65, '[3,8]', FALSE, 1, 'system'),
(6, 'Grammarly', 'grammarly', '✓', '智能语法检查和写作优化工具，帮助提升文档质量。', 'https://www.grammarly.com', 'https://www.grammarly.com/features', '{"free":{"available":true,"limit":"基础检查"},"premium":{"name":"Premium","price":"$12/月","features":["高级语法检查","写作风格建议","学术写作支持"]}}', '免费版基础检查，Premium 版本 $12/月', 1, 'beginner', 4.6, 890, 9200, 72, '[1,3,8]', TRUE, 1, 'system');

-- 图像类工具 (4个)
INSERT INTO `tools` (`id`, `name`, `slug`, `icon`, `description`, `official_url`, `documentation_url`, `pricing`, `pricing_description`, `category_id`, `difficulty`, `rating`, `review_count`, `view_count`, `popularity`, `tags`, `featured`, `status`, `created_by`) VALUES
(7, 'Midjourney', 'midjourney', '🎨', 'AI 图像生成工具，以生成高质量的艺术作品而闻名，支持多种风格和分辨率。', 'https://www.midjourney.com', 'https://docs.midjourney.com', '{"price":"$10/月","features":["无限生成","高分辨率输出","多种风格预设"]}', 'Discord 按生成次数计费，约 $10/月，每月 200 次免费', 2, 'intermediate', 4.8, 1560, 18000, 92, '[4,8,9]', TRUE, 1, 'system'),
(8, 'Stable Diffusion', 'stable-diffusion', '🖼️', '开源的 AI 图像生成模型，可以在本地运行，完全可控，免费使用。', 'https://stability.ai', 'https://platform.stability.ai', '{"free":{"available":true,"limit":"基础功能"},"professional":{"name":"Professional","price":"$20/月","features":["API访问","高优先级","专属资源"]}}', '免费版基础功能，Professional 版本 $20/月', 2, 'intermediate', 4.7, 980, 15000, 85, '[4,7,9]', TRUE, 1, 'system'),
(9, 'DALL-E 3', 'dall-e-3', '🎨', 'OpenAI 开发的文本生成图像模型，理解指令能力强，图像质量高。', 'https://openai.com', 'https://platform.openai.com/docs/guides/images', '{"paid":{"name":"Credit","price":"$0.04/张","features":["高质量图像生成","文本理解"]}}', '按张付费，$0.04/张，每张 1024x1024', 2, 'beginner', 4.5, 720, 14000, 80, '[1,4,7]', TRUE, 1, 'system'),
(10, 'Leonardo AI', 'leonardo-ai', '🎨', 'AI 图像创作平台，提供多种模型和社区分享功能，支持 PSD 导出。', 'https://leonardo.ai', 'https://leonardo.ai', '{"free":{"available":true,"limit":"150 tokens/天"},"pro":{"name":"Pro","price":"10/月","features":["无限tokens","高级模型","团队协作"]}}', '免费版每天 150 tokens，Pro 版本 $10/月', 2, 'beginner', 4.4, 450, 9800, 68, '[4,8,9]', FALSE, 1, 'system');

-- 视频类工具 (3个)
INSERT INTO `tools` (`id`, `name`, `slug`, `icon`, `description`, `official_url`, `documentation_url`, `pricing`, `pricing_description`, `category_id`, `difficulty`, `rating`, `review_count`, `view_count`, `popularity`, `tags`, `featured`, `status`, `created_by`) VALUES
(11, 'Runway', 'runway', '🎬', 'AI 视频生成和编辑工具，提供多种视频特效、背景移除、帧插值等功能。', 'https://runwayml.com', 'https://docs.runwayml.com', '{"free":{"available":true,"limit":"6分钟/月"},"pro":{"name":"Gen-2","price":"$12/月","features":["无限视频生成","高清输出","高级编辑"]}}', '免费版每月 6 分钟，Gen-2 版本 $12/月', 3, 'intermediate', 4.6, 680, 12000, 82, '[5,7,8]', TRUE, 1, 'system'),
(12, 'Pika Labs', 'pika', '🎬', 'AI 视频生成工具，支持文本生成视频和视频风格转换，操作简单。', 'https://pika.art', 'https://pika.art', '{"free":{"available":true,"limit":"30次/天"},"pro":{"name":"Pro","price":"$12/月","features":["无限次数","1080p输出","优先排队"]}}', '免费版每天 30 次，Pro 版本 $12/月', 3, 'beginner', 4.5, 520, 9500, 75, '[5,7,8]', FALSE, 1, 'system'),
(13, 'Sora', 'sora', '🎥', 'OpenAI 开发的文本生成视频模型，支持生成 60 秒高质量视频。', 'https://openai.com', 'https://openai.com', '{"paid":{"name":"Access","price":"按需付费","features":["高质量视频生成","长视频支持"]}}', '目前处于邀请制测试阶段', 3, 'advanced', 4.9, 890, 25000, 98, '[1,5,7]', FALSE, 1, 'system');

-- 音频类工具 (2个)
INSERT INTO `tools` (`id`, `name`, `slug`, `icon`, `description`, `official_url`, `documentation_url`, `pricing`, `pricing_description`, `category_id`, `difficulty`, `rating`, `review_count`, `view_count`, `popularity`, `tags`, `featured`, `status`, `created_by`) VALUES
(14, 'ElevenLabs', 'elevenlabs', '🔊', 'AI 语音合成工具，可以生成极其逼真的人类语音，支持多种语言。', 'https://elevenlabs.io', 'https://docs.elevenlabs.io', '{"free":{"available":true,"limit":"10分钟/月"},"creator":{"name":"Creator","price":"22/月","features":["无限分钟数","150种语音","API访问"]}}', '免费版每月 10 分钟，Creator 版本 $22/月', 4, 'beginner', 4.8, 920, 11000, 82, '[4,7,8]', TRUE, 1, 'system'),
(15, 'Suno', 'suno', '🎵', 'AI 音乐生成工具，可以生成高质量的背景音乐和歌曲。', 'https://suno.ai', 'https://suno.ai', '{"free":{"available":true,"limit":"10首/天"},"unlimited":{"name":"Unlimited","price":"$10/月","features":["无限歌曲","高音质输出"]}}', '免费版每天 10 首，Unlimited 版本 $10/月', 4, 'beginner', 4.5, 680, 8900, 78, '[4,7,8]', FALSE, 1, 'system');

-- 代码类工具 (3个)
INSERT INTO `tools` (`id`, `name`, `slug`, `icon`, `description`, `official_url`, `documentation_url`, `pricing`, `pricing_description`, `category_id`, `difficulty`, `rating`, `review_count`, `view_count`, `popularity`, `tags`, `featured`, `status`, `created_by`) VALUES
(16, 'GitHub Copilot', 'copilot', '💻', 'GitHub 开发的 AI 编程助手，实时提供代码建议和自动补全。', 'https://github.com/features/copilot', 'https://docs.github.com/en/copilot', '{"price":"$10/月","features":["实时代码建议","多语言支持","安全扫描"]}', '个人版 $10/月，企业版额外费用', 5, 'beginner', 4.7, 1250, 18000, 90, '[6,7,8]', TRUE, 1, 'system'),
(17, 'Cursor', 'cursor', '✏️', '基于 VS Code 的 AI 编程编辑器，集成了多种 AI 功能，代码编写效率高。', 'https://cursor.sh', 'https://docs.cursor.sh', '{"free":{"available":true,"limit":"2000次/月"},"pro":{"name":"Pro","price":"$20/月","features":["无限次数","GPT-4支持","Claude支持"]}}', '免费版每月 2000 次对话，Pro 版本 $20/月', 5, 'beginner', 4.6, 890, 14500, 85, '[6,7,8]', TRUE, 1, 'system'),
(18, 'Codeium', 'codeium', '💻', '免费的 AI 编程助手，支持超过 70 种编程语言，提供实时代码建议。', 'https://codeium.com', 'https://codeium.com', '{"free":{"available":true,"features":["无限使用","70+语言支持","无需信用卡"]}}', '完全免费，无需信用卡', 5, 'beginner', 4.5, 620, 12000, 80, '[6,7,9]', TRUE, 1, 'system');

-- 办公类工具 (2个)
INSERT INTO `tools` (`id`, `name`, `slug`, `icon`, `description`, `official_url`, `documentation_url`, `pricing`, `pricing_description`, `category_id`, `difficulty`, `rating`, `review_count`, `view_count`, `popularity`, `tags`, `featured`, `status`, `created_by`) VALUES
(19, 'Microsoft 365 Copilot', 'microsoft-365', '📊', 'Microsoft Office 集成的 AI 助手，帮助撰写文档、分析数据、创建演示文稿。', 'https://www.microsoft.com', 'https://www.microsoft.com/microsoft-365/copilot', '{"price":"$30/月/用户","features":["Word文档AI辅助","Excel数据分析","PPT智能生成"]}', 'Microsoft 365 套餐内包含，单独 $30/月/用户', 6, 'intermediate', 4.5, 450, 9800, 78, '[3,6,8]', FALSE, 1, 'system'),
(20, 'Notion AI', 'notion-ai-office', '📊', 'Notion 的 AI 助手，帮助用户快速整理笔记、生成文档、总结内容。', 'https://www.notion.so', 'https://www.notion.so/product/ai', '{"price":"$10/月","features":["文档AI生成和润色","笔记自动整理","智能搜索"]}', 'Notion 套餐内包含，单独 $10/月', 6, 'beginner', 4.7, 680, 11200, 80, '[3,6,7]', TRUE, 1, 'system');

-- ============================================================================
-- 5. 工具-标签关联数据
-- ============================================================================

INSERT INTO `tool_tag_relations` (`tool_id`, `tag_id`) VALUES
-- ChatGPT 标签
(1, 1), (1, 2), (1, 3), (1, 7), (1, 10),
-- Claude 标签
(2, 1), (2, 2), (2, 3), (2, 7), (2, 10),
-- Midjourney 标签
(7, 4), (7, 8), (7, 9),
-- Stable Diffusion 标签
(8, 4), (8, 7), (8, 9),
-- GitHub Copilot 标签
(16, 6), (16, 7), (16, 9);

-- ============================================================================
-- 6. 徽章数据
-- ============================================================================

INSERT INTO `badges` (`id`, `name`, `slug`, `description`, `icon`, `type`, `condition_type`, `condition_value`, `icon_color`, `background_color`, `status`) VALUES
-- 评测徽章
(1, '首次评测者', 'first-review', '完成第一次工具评测', '📝', 'review', 'review_count', 1, '#FFD700', '#FFF8DC', 1),
(2, '评测达人', 'review-star', '获得5星好评评测', '⭐', 'review', 'review_star', 5, '#FFD700', '#FFF8DC', 1),
(3, '精选评测', 'featured-review', '评测被标记为精选', '💎', 'review', 'featured', 1, '#FFD700', '#FFF8DC', 1),
(4, '深度评测', 'pro-review', '提交50+字详细评测', '🏆', 'review', 'review_count', 50, '#FFD700', '#FFF8DC', 1),
(5, '活跃评测', 'active-review', '每周提交3条评测', '🔥', 'review', 'weekly_review', 3, '#FF4500', '#FFE4E1', 1),
-- 创作徽章
(6, '模板达人', 'template-star', '提交5个模板', '✨', 'contribution', 'template_count', 5, '#87CEEB', '#E0FFFF', 1),
(7, '模板大师', 'template-master', '提交10个模板', '🌟', 'contribution', 'template_count', 10, '#87CEEB', '#E0FFFF', 1),
(8, '模板专家', 'template-expert', '提交20个模板', '💫', 'contribution', 'template_count', 20, '#87CEEB', '#E0FFFF', 1),
-- 社区徽章
(9, '活跃评论', 'active-comment', '发布50条评论', '💬', 'social', 'comment_count', 50, '#32CD32', '#F0FFF0', 1),
(10, '社区达人', 'community-star', '每日登录30天', '📅', 'social', 'login_days', 30, '#32CD32', '#F0FFF0', 1),
(11, '热评', 'hot-comment', '获得100个赞的评论', '👍', 'social', 'like_count', 100, '#FF6347', '#FFE4E1', 1);

-- ============================================================================
-- 7. 系统配置数据
-- ============================================================================

INSERT INTO `system_config` (`key`, `value`, `value_type`, `description`, `category`) VALUES
('tool_review_approve', '1', 'int', '是否需要审核评测', 'review'),
('template_review_approve', '1', 'int', '是否需要审核模板', 'template'),
('enable_anonymous_review', '1', 'boolean', '是否允许匿名评测', 'review'),
('max_review_per_day', '10', 'int', '每日最多评测数', 'review'),
('max_template_per_day', '5', 'int', '每日最多提交模板数', 'template'),
('enable_anonymous_template', '1', 'boolean', '是否允许匿名提交模板', 'template'),
('enable_user_badges', '1', 'boolean', '是否启用徽章系统', 'badge'),
('enable_search_history', '1', 'boolean', '是否记录搜索历史', 'search'),
('enable_tool_views', '1', 'boolean', '是否统计工具浏览量', 'tool');

-- ============================================================================
-- 8. 模拟数据统计（可选）
-- ============================================================================

-- 更新工具的热度值（模拟数据）
UPDATE tools SET popularity = view_count * 0.6 + review_count * 0.3 + (SELECT COALESCE(SUM(likes), 0) FROM tool_reviews tr WHERE tr.tool_id = tools.id) * 0.1 WHERE id BETWEEN 1 AND 20;

-- ============================================================================
-- 9. 验证数据
-- ============================================================================

-- 检查工具类别数量
-- SELECT COUNT(*) as '工具类别数' FROM tool_categories WHERE status = 1;
-- 预期结果: 7

-- 检查工具数量
-- SELECT COUNT(*) as '工具数' FROM tools WHERE status = 1;
-- 预期结果: 20

-- 检查工具标签关联数量
-- SELECT COUNT(*) as '工具标签关联数' FROM tool_tag_relations;
-- 预期结果: 14

-- 检查标签数量
-- SELECT COUNT(*) as '标签数' FROM tool_tags;
-- 预期结果: 10

-- 检查徽章数量
-- SELECT COUNT(*) as '徽章数' FROM badges WHERE status = 1;
-- 预期结果: 11

-- 检查系统配置数量
-- SELECT COUNT(*) as '系统配置数' FROM system_config;
-- 预期结果: 9

-- 查看所有工具（示例）
-- SELECT id, name, category_id, rating, popularity, status FROM tools ORDER BY popularity DESC LIMIT 10;

-- 查看所有工具类别
-- SELECT id, name, icon, featured FROM tool_categories ORDER BY sort_order;

-- 查看工具评分分布
-- SELECT rating, COUNT(*) as count FROM tool_reviews WHERE status = 1 GROUP BY rating ORDER BY rating;

-- ============================================================================
-- 执行完成提示
-- ============================================================================

-- 所有样例数据已加载完成！
-- 您可以运行上面的验证查询来检查数据是否正确。
-- ============================================================================
