-- ======================================================
-- AI热点追踪平台 - 职业风险模块表结构
-- 版本: v1.0
-- 日期: 2026-04-08
-- ======================================================

-- 删除已存在的表（按依赖顺序倒序）
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `profession_management_data`;
DROP TABLE IF EXISTS `profession_market_data`;
DROP TABLE IF EXISTS `profession_transition_advice`;
DROP TABLE IF EXISTS `profession_impact_analysis`;
DROP TABLE IF EXISTS `professions`;
DROP TABLE IF EXISTS `profession_categories`;
SET FOREIGN_KEY_CHECKS = 1;

USE `hot_ai`;

-- ============================================
-- 1. 职业分类表
-- ============================================
CREATE TABLE `profession_categories` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '分类ID',
    `name` varchar(50) NOT NULL COMMENT '分类名称',
    `description` varchar(255) DEFAULT NULL COMMENT '分类描述',
    `sort_order` int DEFAULT '0' COMMENT '排序权重',
    `status` tinyint DEFAULT '1' COMMENT '状态: 0-禁用, 1-启用',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='职业分类表';

-- ============================================
-- 2. 职业信息表（核心表）
-- ============================================
CREATE TABLE `professions` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '职业ID',
    `name` varchar(100) NOT NULL COMMENT '职业名称',
    `slug` varchar(100) NOT NULL COMMENT 'URL友好标识',
    `icon` varchar(20) DEFAULT NULL COMMENT 'Emoji图标',
    `category_id` int DEFAULT NULL COMMENT '分类ID',
    `description` text COMMENT '职业描述',
    `risk_level` enum('extreme','high','medium','low') NOT NULL DEFAULT 'medium' COMMENT '风险等级: extreme-极高风险, high-高风险, medium-中风险, low-低风险',
    `risk_score` int NOT NULL DEFAULT '50' COMMENT '风险指数 (0-100)',
    `automation_rate` int NOT NULL DEFAULT '50' COMMENT '自动化率 (0-100)',
    `status` tinyint DEFAULT '1' COMMENT '状态: 0-待审核, 1-已发布, 2-已删除',
    `sort_order` int DEFAULT '0' COMMENT '排序权重',
    `published_at` datetime DEFAULT NULL COMMENT '发布时间',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='职业信息表';

-- ============================================
-- 3. 职业影响分析表
-- ============================================
CREATE TABLE `profession_impact_analysis` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '影响分析ID',
    `profession_id` int NOT NULL COMMENT '职业ID',
    `affected_tasks` json DEFAULT NULL COMMENT '可被AI替代的任务列表',
    `safe_tasks` json DEFAULT NULL COMMENT '难以被AI替代的任务列表',
    `safe_skills` json DEFAULT NULL COMMENT '不可替代技能列表',
    `impact_timeline` json DEFAULT NULL COMMENT '影响时间线',
    `impact_summary` text COMMENT '影响分析总结',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='职业影响分析表';

-- ============================================
-- 4. 职业转型建议表
-- ============================================
CREATE TABLE `profession_transition_advice` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '转型建议ID',
    `profession_id` int NOT NULL COMMENT '职业ID',
    `transition_paths` json DEFAULT NULL COMMENT '转型方向建议列表',
    `recommended_skills` json DEFAULT NULL COMMENT '推荐学习技能',
    `recommended_tools` json DEFAULT NULL COMMENT '推荐工具ID列表',
    `recommended_paths` json DEFAULT NULL COMMENT '推荐学习路径ID列表',
    `related_articles` json DEFAULT NULL COMMENT '相关文章ID列表',
    `advice_summary` text COMMENT '转型建议总结',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='职业转型建议表';

-- ============================================
-- 5. 职业市场数据表
-- ============================================
CREATE TABLE `profession_market_data` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '市场数据ID',
    `profession_id` int NOT NULL COMMENT '职业ID',
    `market_trend` enum('growing','stable','declining') DEFAULT 'stable' COMMENT '市场需求趋势: growing-增长, stable-稳定, declining-下降',
    `market_trend_description` text COMMENT '市场趋势说明',
    `salary_impact` enum('positive','neutral','negative') DEFAULT 'neutral' COMMENT '薪资影响: positive-正面, neutral-中性, negative-负面',
    `salary_change_rate` decimal(5,2) DEFAULT NULL COMMENT '薪资变化率(%)',
    `avg_salary` decimal(10,2) DEFAULT NULL COMMENT '平均薪资(元/月)',
    `job_demand_trend` varchar(50) DEFAULT NULL COMMENT '岗位需求趋势描述',
    `supply_demand_ratio` decimal(3,2) DEFAULT NULL COMMENT '供需比例',
    `data_source` varchar(255) DEFAULT NULL COMMENT '数据来源',
    `data_update_date` date DEFAULT NULL COMMENT '数据更新日期',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='职业市场数据表';

-- ============================================
-- 6. 职业管理数据表
-- ============================================
CREATE TABLE `profession_management_data` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '管理数据ID',
    `profession_id` int NOT NULL COMMENT '职业ID',
    `view_count` int DEFAULT '0' COMMENT '查看次数',
    `search_count` int DEFAULT '0' COMMENT '搜索次数',
    `favorite_count` int DEFAULT '0' COMMENT '收藏次数',
    `meta_title` varchar(255) DEFAULT NULL COMMENT 'SEO标题',
    `meta_description` varchar(500) DEFAULT NULL COMMENT 'SEO描述',
    `meta_keywords` varchar(500) DEFAULT NULL COMMENT 'SEO关键词',
    `reviewer_id` varchar(36) DEFAULT NULL COMMENT '审核人ID',
    `reviewed_at` datetime DEFAULT NULL COMMENT '审核时间',
    `review_notes` text COMMENT '审核备注',
    `is_featured` tinyint DEFAULT '0' COMMENT '是否精选: 0-否, 1-是',
    `featured_at` datetime DEFAULT NULL COMMENT '精选时间',
    `tags` json DEFAULT NULL COMMENT '运营标签',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='职业管理数据表';
