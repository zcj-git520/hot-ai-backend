-- AI 热点追踪平台 - 爬虫服务数据库表结构
-- 版本：v1.0
-- 创建日期：2026-04-02

USE hot_ai;

-- ============================================
-- 抓取源表
-- ============================================
CREATE TABLE IF NOT EXISTS crawler_sources (
    id VARCHAR(36) PRIMARY KEY COMMENT '抓取源 ID',
    name VARCHAR(100) NOT NULL COMMENT '抓取源名称',
    url VARCHAR(500) NOT NULL COMMENT '抓取目标 URL',
    source_type ENUM('rss', 'sitemap', 'api', 'html') NOT NULL DEFAULT 'html' COMMENT '抓取源类型',
    category VARCHAR(50) COMMENT '内容分类（如：AI 技术、AI 产品、AI 政策等）',
    fetch_method ENUM('get', 'post') NOT NULL DEFAULT 'get' COMMENT '请求方法',
    request_headers JSON COMMENT '自定义请求头',
    request_body TEXT COMMENT '请求体（POST 时使用）',
    
    -- 解析配置
    parser_type ENUM('json', 'xpath', 'css', 'regex', 'html') NOT NULL DEFAULT 'html' COMMENT '解析器类型',
    parse_rules JSON COMMENT '解析规则配置',
    
    -- 调度配置
    fetch_interval INT DEFAULT 300 COMMENT '抓取间隔（秒），默认 5 分钟',
    timeout INT DEFAULT 30 COMMENT '超时时间（秒）',
    max_retries INT DEFAULT 3 COMMENT '最大重试次数',
    priority INT DEFAULT 1 COMMENT '优先级（1-10，越高越优先）',
    
    -- 状态管理
    status ENUM('active', 'inactive', 'error', 'paused') NOT NULL DEFAULT 'active' COMMENT '抓取源状态',
    last_fetch_at DATETIME COMMENT '上次抓取时间',
    next_fetch_at DATETIME COMMENT '下次计划抓取时间',
    last_status_code INT COMMENT '上次抓取 HTTP 状态码',
    last_error_message TEXT COMMENT '上次错误信息',
    consecutive_failures INT DEFAULT 0 COMMENT '连续失败次数',
    
    -- 统计信息
    total_fetches INT DEFAULT 0 COMMENT '累计抓取次数',
    successful_fetches INT DEFAULT 0 COMMENT '成功抓取次数',
    failed_fetches INT DEFAULT 0 COMMENT '失败抓取次数',
    total_items INT DEFAULT 0 COMMENT '累计抓取条目数',
    
    -- 元数据
    description TEXT COMMENT '抓取源描述',
    tags JSON COMMENT '标签',
    created_by VARCHAR(36) COMMENT '创建人',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_status (status),
    INDEX idx_category (category),
    INDEX idx_source_type (source_type),
    INDEX idx_priority (priority DESC),
    INDEX idx_next_fetch (next_fetch_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='抓取源配置表';

-- ============================================
-- 抓取日志表
-- ============================================
CREATE TABLE IF NOT EXISTS crawler_fetch_logs (
    id VARCHAR(36) PRIMARY KEY COMMENT '日志 ID',
    source_id VARCHAR(36) NOT NULL COMMENT '抓取源 ID',
    fetch_started_at DATETIME NOT NULL COMMENT '开始抓取时间',
    fetch_completed_at DATETIME COMMENT '完成抓取时间',
    duration_ms INT COMMENT '耗时（毫秒）',
    status_code INT COMMENT 'HTTP 状态码',
    response_size INT COMMENT '响应大小（字节）',
    items_fetched INT DEFAULT 0 COMMENT '抓取条目数',
    status ENUM('success', 'failed', 'timeout', 'skipped') NOT NULL COMMENT '抓取状态',
    error_message TEXT COMMENT '错误信息',
    raw_response LONGTEXT COMMENT '原始响应内容（可选，用于调试）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    FOREIGN KEY (source_id) REFERENCES crawler_sources(id) ON DELETE CASCADE,
    INDEX idx_source_id (source_id),
    INDEX idx_status (status),
    INDEX idx_fetch_time (fetch_started_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='抓取日志表';

-- ============================================
-- 初始化示例数据
-- ============================================

-- 插入示例抓取源（以国外主流 AI 资讯源为主，辅以国内优质源）
INSERT INTO crawler_sources (id, name, url, source_type, category, fetch_interval, priority, status, description) VALUES
-- ========== 国外主流 AI 资讯源 ==========
('src_001', 'OpenAI Blog', 'https://openai.com/blog/', 'rss', 'AI 技术', 1800, 10, 'active', 'OpenAI 官方博客 - GPT、DALL-E 等前沿研究'),
('src_002', 'Google AI Blog', 'https://ai.googleblog.com/', 'rss', 'AI 技术', 1800, 10, 'active', 'Google AI 官方博客 - Transformer、Gemini 等'),
('src_003', 'DeepMind Blog', 'https://deepmind.google/discover/blog/', 'rss', 'AI 技术', 1800, 10, 'active', 'DeepMind 博客 - AlphaGo、AlphaFold 等突破性研究'),
('src_004', 'Hugging Face Blog', 'https://huggingface.co/blog', 'rss', 'AI 技术', 1800, 9, 'active', 'Hugging Face 技术博客 - NLP、开源模型'),
('src_005', 'Anthropic Blog', 'https://www.anthropic.com/news', 'rss', 'AI 技术', 1800, 9, 'active', 'Anthropic 官方博客 - Claude、AI 安全研究'),
('src_006', 'MIT Technology Review - AI', 'https://www.technologyreview.com/topic/artificial-intelligence/', 'html', 'AI 产业', 3600, 8, 'active', '麻省理工科技评论 AI 板块 - 深度报道'),
('src_007', 'VentureBeat AI', 'https://venturebeat.com/category/ai/', 'html', 'AI 产业', 1800, 8, 'active', 'VentureBeat AI - AI 产业新闻与商业应用'),
('src_008', 'The Verge - AI', 'https://www.theverge.com/ai-artificial-intelligence', 'html', 'AI 产业', 3600, 7, 'active', 'The Verge AI 板块 - 科技媒体深度报道'),
('src_009', 'Wired - AI', 'https://www.wired.com/tag/artificial-intelligence/', 'html', 'AI 产业', 3600, 7, 'active', 'Wired AI 相关报道 - 技术与社会影响'),
('src_010', 'Ars Technica - AI', 'https://arstechnica.com/ai/', 'html', 'AI 产业', 3600, 7, 'active', 'Ars Technica AI 新闻 - 技术分析'),
('src_011', 'TechCrunch - AI', 'https://techcrunch.com/category/artificial-intelligence/', 'html', 'AI 创业', 1800, 8, 'active', 'TechCrunch AI 创业与投资动态'),
('src_012', 'Stanford HAI', 'https://hai.stanford.edu/news', 'rss', 'AI 研究', 3600, 8, 'active', '斯坦福大学以人为本 AI 研究院'),
('src_013', 'Berkeley AI Research', 'https://bair.berkeley.edu/blog/', 'rss', 'AI 研究', 3600, 8, 'active', '伯克利 AI 研究实验室博客'),
('src_014', 'CMU ML News', 'https://www.ml.cmu.edu/news', 'rss', 'AI 研究', 3600, 7, 'active', '卡内基梅隆大学机器学习新闻'),
('src_015', 'AI2 Blog', 'https://allenai.org/blog', 'rss', 'AI 研究', 3600, 7, 'active', '艾伦人工智能研究所官方博客'),

-- ========== 国内主流 AI 资讯源 ==========
('src_016', '机器之心', 'https://www.jiqizhixin.com/', 'html', 'AI 资讯', 600, 9, 'active', '机器之心 - 国内领先 AI 媒体，专注前沿技术'),
('src_017', '量子位', 'https://www.qbitai.com/', 'html', 'AI 资讯', 600, 9, 'active', '量子位 - AI 前沿动态与产业观察'),
('src_018', '新智元', 'https://www.newsmartbasic.com/', 'html', 'AI 资讯', 600, 8, 'active', '新智元 - AI 产业生态新媒体'),
('src_019', 'AI 科技评论', 'https://www.aitechtalk.com/', 'html', 'AI 技术', 1200, 8, 'active', 'AI 科技评论 - 学术界与产业界桥梁'),
('src_020', '雷锋网', 'https://www.leiphone.com/category/ai', 'html', 'AI 产业', 1200, 7, 'active', '雷锋网 AI 板块 - 智能硬件与 AI 结合')
ON DUPLICATE KEY UPDATE name=name;
