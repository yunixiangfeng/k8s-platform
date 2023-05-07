CREATE TABLE `workflow` ( 
	`id` int NOT NULL AUTO_INCREMENT,
	`name` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
	`namespace` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
	`replicas` int DEFAULT NULL,
	`deployment` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
	`service` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
	`ingress` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
	`type` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL,
	`created_at` datetime DEFAULT NULL,
	`updated_at` datetime DEFAULT NULL,
	`deleted_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE KEY `name` (`name`)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
