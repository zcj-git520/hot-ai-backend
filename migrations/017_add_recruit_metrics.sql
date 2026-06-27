USE `hot_ai`;

ALTER TABLE `profession_market_data`
  ADD COLUMN `market_confidence_score` TINYINT UNSIGNED DEFAULT NULL AFTER `data_update_date`,
  ADD COLUMN `last_metrics_date` DATE DEFAULT NULL AFTER `market_confidence_score`,
  ADD COLUMN `metric_window` VARCHAR(20) DEFAULT '30d' AFTER `last_metrics_date`,
  ADD COLUMN `data_freshness` ENUM('fresh','stale','missing') NOT NULL DEFAULT 'missing' AFTER `metric_window`,
  ADD INDEX `idx_data_freshness` (`data_freshness`);
