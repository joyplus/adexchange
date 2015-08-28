ALTER TABLE `pmp_demand_adspace` ADD `adspace_type` int(1) NOT NULL AFTER `app_id`;
ALTER TABLE `pmp_request_log` ADD `did` VARCHAR(50) NOT NULL AFTER `bid`;
ALTER TABLE `pmp_tracking_log` ADD `did` VARCHAR(50) NOT NULL AFTER `bid`;
ALTER TABLE `pmp_demand_response_log` ADD `did` VARCHAR(50) NOT NULL AFTER `bid`;

/*20150821*/
ALTER TABLE `pmp_adspace` ADD `tpl_name` VARCHAR(20) NOT NULL AFTER `creative_type`;

/*20150825*/
CREATE TABLE `pmp_campaign_creative` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `campaign_id` int(11) NOT NULL,
  `name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `width` int(11) DEFAULT NULL,
  `height` int(11) DEFAULT NULL,
  `creative_url` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `creative_status` int(11) DEFAULT NULL COMMENT '0：暂停1： 运行',
  `landing_url` varchar(500) COLLATE utf8_unicode_ci DEFAULT NULL,
  `imp_tracking_url` varchar(1000) COLLATE utf8_unicode_ci DEFAULT NULL,
  `clk_tracking_url` varchar(1000) COLLATE utf8_unicode_ci NOT NULL,
  `display_title` varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  `display_text` varchar(1000) COLLATE utf8_unicode_ci DEFAULT NULL，
  PRIMARY KEY (`id`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 COLLATE=utf8_unicode_ci;

ALTER TABLE `pmp_adspace` ADD `forever_flg` INT(2) NOT NULL AFTER `media_id`;

/*20150828*/
ALTER TABLE `pmp_campaign`
  DROP `width`,
  DROP `height`,
  DROP `creative_url`;