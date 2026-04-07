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
('提示词工程', 3),
('Gemini', 1),
('Claude', 1),
('代码生成', 1),
('图像生成', 1),
('视频生成', 1),
('AI伦理', 1),
('自动驾驶', 1),
('机器人', 1),
('医疗AI', 1),
('金融科技', 1),
('文案写作', 2),
('翻译', 2),
('数据分析', 2),
('产品经理', 2),
('市场营销', 2),
('教育行业', 2),
('医疗行业', 2),
('法律行业', 2),
('财务分析', 2),
('TensorFlow', 3),
('PyTorch', 3),
('机器学习', 3),
('神经网络', 3),
('强化学习', 3),
('数据科学', 3),
('SQL', 3),
('JavaScript', 3),
('Prompt Engineering', 3),
('微调', 3),
('模型部署', 3);

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
),
(
  'Google Gemini Ultra 2.0 发布：性能超越 GPT-5',
  'Google 今日发布 Gemini Ultra 2.0，在多模态理解和推理能力上实现重大突破，特别是在数学和科学领域的表现尤为突出...',
  '507f1f77bcf86cd799439017',
  2,
  '科技前沿',
  1,
  '2026-04-01 09:15:00',
  1
),
(
  'Anthropic Claude 4 上线：更安全的 AI 助手',
  'Anthropic 发布 Claude 4，强调安全性和可靠性，在有害内容过滤和价值观对齐方面取得显著进步...',
  '507f1f77bcf86cd799439018',
  1,
  '安全AI',
  1,
  '2026-04-02 14:20:00',
  1
),
(
  'AI 视频生成新突破：Runway Gen-3 震撼登场',
  'Runway ML 发布 Gen-3 视频生成模型，支持高分辨率长视频生成，为电影制作和内容创作带来革命性变革...',
  '507f1f77bcf86cd799439019',
  3,
  '视觉科技',
  4,
  '2026-04-03 11:45:00',
  1
),
(
  'AI 在医疗诊断中的最新应用进展',
  '深度学习模型在医学影像分析、疾病预测和药物研发等领域取得重大突破，本文综述了 AI 在医疗领域的最新应用案例...',
  '507f1f77bcf86cd799439020',
  5,
  '医疗科技',
  1,
  '2026-04-04 10:30:00',
  1
),
(
  '金融行业如何利用 AI 进行风险评估？',
  '从信用评分到欺诈检测，AI 正在彻底改变金融风控体系。本文介绍了机器学习在金融风险评估中的实际应用...',
  '507f1f77bcf86cd799439021',
  6,
  '金融科技',
  2,
  '2026-04-05 16:00:00',
  1
),
(
  'TensorFlow 2.15 发布：性能优化与新特性',
  'Google 发布 TensorFlow 2.15，重点优化了大模型训练效率，并引入了新的分布式训练工具链...',
  '507f1f77bcf86cd799439022',
  2,
  '开发者社区',
  3,
  '2026-04-06 13:10:00',
  1
),
(
  '从零开始：如何微调自己的大语言模型？',
  '本文详细介绍了微调大语言模型的完整流程，包括数据准备、模型选择、训练技巧和评估方法...',
  '507f1f77bcf86cd799439023',
  4,
  'AI教程',
  3,
  '2026-04-07 09:30:00',
  1
),
(
  'AI 伦理与治理：全球监管趋势分析',
  '随着 AI 技术快速发展，各国纷纷出台相关法规。本文分析了欧盟 AI 法案、美国 AI 治理框架等主要监管政策...',
  '507f1f77bcf86cd799439024',
  1,
  '政策观察',
  2,
  '2026-04-08 15:20:00',
  1
),
(
  '自动驾驶 Level 4 商业化落地面临哪些挑战？',
  '尽管技术不断进步，但自动驾驶的商业化仍然面临法规、安全和成本等多重挑战。本文深入分析了当前瓶颈...',
  '507f1f77bcf86cd799439025',
  3,
  '汽车科技',
  1,
  '2026-04-09 11:00:00',
  1
),
(
  'AI 内容审核：如何平衡安全与言论自由？',
  '随着 AI 生成内容泛滥，内容审核成为平台面临的重大挑战。本文探讨了 AI 审核技术的现状与伦理困境...',
  '507f1f77bcf86cd799439026',
  4,
  '科技伦理',
  2,
  '2026-04-10 14:00:00',
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
(6, 1567, 28, 178),
(7, 3120, 45, 420),
(8, 2789, 38, 356),
(9, 1890, 29, 245),
(10, 2456, 42, 321),
(11, 1987, 36, 278),
(12, 2234, 41, 302),
(13, 2678, 48, 387),
(14, 1892, 32, 256),
(15, 2345, 39, 298),
(16, 2100, 35, 274);

-- ============================================
-- 7. 文章-标签关联（补充文章）
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

-- 文章7: Google Gemini Ultra 2.0
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(7, 16),  -- Gemini
(7, 11),  -- 大模型
(7, 12),  -- 自然语言处理
(7, 13);  -- 计算机视觉

-- 文章8: Anthropic Claude 4
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(8, 17),  -- Claude
(8, 11),  -- 大模型
(8, 21),  -- AI伦理
(8, 1);   -- GPT-5 (相关)

-- 文章9: AI视频生成新突破
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(9, 20),  -- 视频生成
(9, 19),  -- 图像生成
(9, 14),  -- AIGC
(9, 13);  -- 计算机视觉

-- 文章10: AI医疗诊断
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(10, 24),  -- 医疗AI
(10, 10),  -- 深度学习
(10, 13),  -- 计算机视觉
(10, 38);  -- 神经网络

-- 文章11: 金融AI风险评估
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(11, 25),  -- 金融科技
(11, 28),  -- 数据分析
(11, 37),  -- 机器学习
(11, 40);  -- 数据科学

-- 文章12: TensorFlow 2.15
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(12, 35),  -- TensorFlow
(12, 10),  -- 深度学习
(12, 45),  -- 模型部署
(12, 37);  -- 机器学习

-- 文章13: 微调大语言模型
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(13, 44),  -- 微调
(13, 11),  -- 大模型
(13, 10),  -- 深度学习
(13, 43);  -- Prompt Engineering

-- 文章14: AI伦理与治理
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(14, 21),  -- AI伦理
(14, 3),   -- 职业影响
(14, 6);   -- 职业发展

-- 文章15: 自动驾驶
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(15, 22),  -- 自动驾驶
(15, 23),  -- 机器人
(15, 13),  -- 计算机视觉
(15, 10);  -- 深度学习

-- 文章16: AI内容审核
INSERT INTO `article_tag_relation` (`article_id`, `tag_id`) VALUES
(16, 21),  -- AI伦理
(16, 12),  -- 自然语言处理
(16, 13),  -- 计算机视觉
(16, 14);  -- AIGC

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
