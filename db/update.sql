ALTER TABLE `pmp_request_log` ADD `did` VARCHAR(50) NOT NULL AFTER `bid`;
ALTER TABLE `pmp_tracking_log` ADD `did` VARCHAR(50) NOT NULL AFTER `bid`;
ALTER TABLE `pmp_demand_response_log` ADD `did` VARCHAR(50) NOT NULL AFTER `bid`;

/*20150821*/
ALTER TABLE `pmp_adspace` ADD `tpl_name` VARCHAR(20) NOT NULL AFTER `creative_type`;