INSERT INTO `pmp`.`pmp_demand_platform_desk` (`id`, `name`, `request_url_template`, `timeout`, `invoke_func_name`) VALUES ('4', 'baidu', 'http://220.181.163.105/api', '100', 'invokeBD');

-- http://mobads.baidu.com/api

INSERT INTO `pmp`.`pmp_demand_adspace` (`id`, `name`, `demand_adspace_key`, `demand_id`) VALUES ('8', 'Baidu testing demand side', 'L000015a', '4');

INSERT INTO `pmp`.`pmp_daily_allocation` (`id`, `ad_date`, `demand_adspace_id`, `imp`, `pmp_adspace_id`) VALUES ('21', '2015-07-11', '8', '200', '3');

INSERT INTO `pmp`.`pmp_adspace_matrix` (`id`, `pmp_adspace_id`, `demand_id`, `demand_adspace_id`, `priority`) VALUES ('6', '3', '4', '8', '4');
