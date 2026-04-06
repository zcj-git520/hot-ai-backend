-- ======================================================
-- AI热点追踪平台 - 文章模块样例数据
-- 版本: v1.0
-- 日期: 2026-04-06
-- ======================================================

USE hot_ai;

-- ============================================
-- 1. 分类数据
-- ============================================
INSERT INTO `categories` (`name`, `code`, `color`, `icon`, `sort_order`) VALUES
('动态', 'news', '#3B82F6', 'news', 1),
('职业', 'impact', '#F97316', 'work', 2),
('学习', 'learn', '#10B981', 'school', 3),
('工具', 'tool', '#8B5CF6', 'build', 4);

-- ============================================
-- 2. 来源媒体数据
-- ============================================
INSERT INTO `sources` (`name`, `domain`, `logo_url`, `description`, `reliability_score`) VALUES
('机器之心', 'jiqizhixin.com', 'https://cdn.jiqizhixin.com/logo.png', '专注人工智能领域的前沿科技报道', 9),
('量子位', 'qbitai.com', 'https://cdn.qbitai.com/logo.png', 'AI与前沿科技媒体', 9),
('知乎 AI', 'zhihu.com', 'https://cdn.zhihu.com/logo.png', '知乎人工智能话题精选', 8),
('36氪', '36kr.com', 'https://cdn.36kr.com/logo.png', '创业投资新媒体', 8),
('AI科技评论', 'aitecher.com', 'https://cdn.aitecher.com/logo.png', '人工智能学术与产业观察', 9),
('雷锋网', 'leiphone.com', 'https://cdn.leiphone.com/logo.png', '智能硬件第一媒体', 7);

-- ============================================
-- 3. 标签数据
-- ============================================
INSERT INTO `tags` (`name`, `type`) VALUES
('GPT-5', 1),
('OpenAI', 1),
('职业影响', 2),
('设计师', 2),
('转型', 2),
('职业发展', 2),
('学习路线', 3),
('零基础', 3),
('Python', 3),
('深度学习', 3),
('大模型', 1),
('自然语言处理', 1),
('计算机视觉', 1),
('AIGC', 1),
('提示词工程', 3);

-- ============================================
-- 4. 文章数据（MongoDB ObjectId为模拟数据）
-- ============================================
INSERT INTO `articles` (`title`, `summary`, `content_mongo_id`, `source_id`, `author`, `category_id`, `published_at`, `status`) VALUES
(
  'GPT-5 发布：AI 能力再次飞跃，这些职业将受到影响',
  'OpenAI 今日发布 GPT-5，在多个基准测试中取得突破性进展。专家分析，文案、设计、客服等岗位将面临更大挑战，但同时也催生了新的职业机会...',
  '507f1f77bcf86cd799439011',
  1,
  'AI观察',
  1,
  '2026-03-30 10:00:00',
  1
),
(
  '设计师如何应对 AI 冲击？这份转型指南请收好',
  'AI 绘图工具的崛起让设计师群体倍感压力。本文采访了 5 位成功转型的设计师，分享他们的经验和建议，包括转向UI/UX设计、品牌策略等领域...',
  '507f1f77bcf86cd799439012',
  3,
  '设计前沿',
  2,
  '2026-03-29 14:30:00',
  1
),
(
  '零基础学习 AI 的完整路线图（2026 年版）',
  '从 Python 基础到深度学习，从理论到实践，这份学习路线图为初学者提供了清晰的学习路径，包含免费资源和实战项目推荐...',
  '507f1f77bcf86cd799439013',
  2,
  'AI教育',
  3,
  '2026-03-28 09:00:00',
  1
),
(
  'Midjourney V7 震撼发布：图像生成质量再创新高',
  'Midjourney 团队今日发布 V7 版本，在图像细节、文字渲染、多主体构图等方面均有显著提升，进一步缩小了与真实照片的差距...',
  '507f1f77bcf86cd799439014',
  1,
  '科技日报',
  4,
  '2026-03-27 16:00:00',
  1
),
(
  '程序员必看：AI 辅助编程的最佳实践',
  'GitHub Copilot、Cursor 等 AI 编程工具已经普及，但如何高效使用它们？本文总结了 10 条最佳实践，帮助开发者提升工作效率...',
  '507f1f77bcf86cd799439015',
  5,
  '代码工匠',
  3,
  '2026-03-26 11:00:00',
  1
),
(
  'AI 产品经理的核心竞争力是什么？',
  '随着 AI 产品越来越多，市场对 AI 产品经理的需求也在增长。本文分析了优秀 AI 产品经理需要具备的 5 项核心能力...',
  '507f1f77bcf86cd799439016',
  4,
  '产品思维',
  2,
  '2026-03-25 13:30:00',
  1
);

-- ============================================
-- 5. 文章统计数据
-- ============================================
INSERT INTO `article_stats` (`article_id`, `view_count`, `comment_count`, `like_count`) VALUES
(1, 1234, 23, 156),
(2, 2456, 45, 289),
(3, 3789, 67, 423),
(4, 1890, 34, 201),
(5, 2100, 56, 312),
(6, 1567, 28, 178);

-- ============================================
-- 6. 文章-标签关联
-- ============================================
-- 文章1: GPT-5 发布
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(1, 1),  -- GPT-5
(1, 2),  -- OpenAI
(1, 3),  -- 职业影响
(1, 11); -- 大模型

-- 文章2: 设计师转型
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(2, 4),  -- 设计师
(2, 5),  -- 转型
(2, 6),  -- 职业发展
(2, 14); -- AIGC

-- 文章3: 学习路线
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(3, 7),  -- 学习路线
(3, 8),  -- 零基础
(3, 9),  -- Python
(3, 10); -- 深度学习

-- 文章4: Midjourney
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(4, 13), -- 计算机视觉
(4, 14); -- AIGC

-- 文章5: AI编程
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(5, 9),  -- Python
(5, 15); -- 提示词工程

-- 文章6: AI产品经理
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(6, 3),  -- 职业影响
(6, 6);  -- 职业发展

-- ============================================
-- 验证数据
-- ============================================
SELECT '✅ 文章模块数据初始化完成!' AS message;
SELECT COUNT(*) AS categories_count FROM categories;
SELECT COUNT(*) AS sources_count FROM sources;
SELECT COUNT(*) AS tags_count FROM tags;
SELECT COUNT(*) AS articles_count FROM articles;
SELECT COUNT(*) AS article_stats_count FROM article_stats;
SELECT COUNT(*) AS article_tag_relations_count FROM article_tag_relation;
