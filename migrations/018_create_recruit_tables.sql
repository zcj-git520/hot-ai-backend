USE `hot_ai`;

CREATE TABLE `recruit_raw_jobs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `platform` ENUM('boss','zhilian','liepin') NOT NULL,
  `platform_job_id` VARCHAR(128) NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `company` VARCHAR(255) DEFAULT NULL,
  `city` VARCHAR(50) DEFAULT NULL,
  `salary_min` INT DEFAULT NULL,
  `salary_max` INT DEFAULT NULL,
  `description` TEXT,
  `skills` JSON DEFAULT NULL,
  `industry` VARCHAR(100) DEFAULT NULL,
  `url` VARCHAR(512) DEFAULT NULL,
  `crawled_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_platform_job` (`platform`, `platform_job_id`),
  KEY `idx_crawled_at` (`crawled_at`),
  KEY `idx_city` (`city`),
  KEY `idx_industry` (`industry`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `recruit_normalized_jobs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `raw_job_id` BIGINT UNSIGNED NOT NULL,
  `profession_id` INT UNSIGNED DEFAULT NULL,
  `match_method` ENUM('llm','keyword','manual') NOT NULL,
  `match_confidence` DECIMAL(3,2) DEFAULT NULL,
  `ai_keywords_count` INT DEFAULT 0,
  `ai_keywords_total` INT DEFAULT 0,
  `ai_keyword_hits` JSON DEFAULT NULL,
  `normalized_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_raw_job_id` (`raw_job_id`),
  KEY `idx_profession_time` (`profession_id`, `normalized_at`),
  CONSTRAINT `fk_normalized_raw` FOREIGN KEY (`raw_job_id`) REFERENCES `recruit_raw_jobs` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `recruit_daily_metrics` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `profession_id` INT UNSIGNED NOT NULL,
  `metric_date` DATE NOT NULL,
  `job_count` INT NOT NULL DEFAULT 0,
  `job_count_prev_30d` INT NOT NULL DEFAULT 0,
  `salary_median` DECIMAL(10,2) DEFAULT NULL,
  `salary_median_prev_90d` DECIMAL(10,2) DEFAULT NULL,
  `ai_penetration_rate` DECIMAL(4,2) DEFAULT NULL,
  `geo_distribution` JSON DEFAULT NULL,
  `industry_distribution` JSON DEFAULT NULL,
  `sample_size` INT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_profession_date` (`profession_id`, `metric_date`),
  KEY `idx_metric_date` (`metric_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `recruit_keywords` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `profession_id` INT UNSIGNED NOT NULL,
  `keyword` VARCHAR(100) NOT NULL,
  `weight` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_profession_keyword` (`profession_id`, `keyword`),
  KEY `idx_profession` (`profession_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `recruit_crawl_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `platform` ENUM('boss','zhilian','liepin') NOT NULL,
  `started_at` DATETIME NOT NULL,
  `finished_at` DATETIME DEFAULT NULL,
  `status` ENUM('success','partial','failed') NOT NULL,
  `jobs_fetched` INT NOT NULL DEFAULT 0,
  `jobs_new` INT NOT NULL DEFAULT 0,
  `error_message` TEXT,
  PRIMARY KEY (`id`),
  KEY `idx_platform_time` (`platform`, `started_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
