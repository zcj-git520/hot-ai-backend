USE `hot_ai`;

CREATE TABLE IF NOT EXISTS `recruit_config` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `config_key` VARCHAR(50) NOT NULL,
  `config_value` DECIMAL(10,3) NOT NULL,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `recruit_config` (`config_key`, `config_value`) VALUES
  ('weight.demand_decay', 0.400),
  ('weight.salary_drop', 0.300),
  ('weight.ai_penetration', 0.200),
  ('weight.distribution_concentration', 0.100),
  ('threshold.stale_sample_size', 30.000),
  ('threshold.fresh_min_days', 7.000);
