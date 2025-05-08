CREATE TABLE IF NOT EXISTS `apps` (
  `id` int NOT NULL AUTO_INCREMENT,
  `app_key` varchar(45) NOT NULL,
  `app_secret` varchar(45) NOT NULL,
  `app_secure_key` varchar(45) NOT NULL,
  `app_status` tinyint DEFAULT '0',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_type` tinyint DEFAULT '0',
  `app_name` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `appexts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `app_key` varchar(50) DEFAULT NULL,
  `app_item_key` varchar(50) DEFAULT NULL,
  `app_item_value` varchar(2048) DEFAULT NULL,
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `IDX_APPKEY_APPITEMKEY` (`app_key`,`app_item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `fileconfs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `app_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `channel` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `conf` json DEFAULT NULL,
  `enable` tinyint(1) DEFAULT '0',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `app_key` (`app_key`,`channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `friendrels` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `friend_id` varchar(32) DEFAULT NULL,
  `order_tag` varchar(20) NULL DEFAULT '',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_friend` (`app_key`,`user_id`,`friend_id`),
  KEY `idx_order` (`app_key`, `user_id`, `order_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `groupadmins` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group_id` varchar(64) DEFAULT NULL,
  `admin_id` varchar(64) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_admin` (`app_key`,`group_id`,`admin_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `groupinfos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group_id` varchar(64) DEFAULT NULL,
  `group_name` varchar(64) DEFAULT NULL,
  `group_portrait` varchar(200) DEFAULT NULL,
  `creator_id` varchar(64) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  `is_mute` tinyint DEFAULT '0',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_groupid` (`app_key`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `groupinfoexts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group_id` varchar(32) DEFAULT NULL,
  `item_key` varchar(50) DEFAULT NULL,
  `item_value` varchar(100) DEFAULT NULL,
  `item_type` tinyint DEFAULT '0',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_groupid` (`app_key`,`group_id`,`item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `groupmembers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group_id` varchar(64) DEFAULT NULL,
  `member_id` varchar(64) DEFAULT NULL,
  `member_type` tinyint DEFAULT '0',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(45) DEFAULT NULL,
  `is_mute` tinyint DEFAULT '0',
  `is_allow` tinyint DEFAULT '0',
  `mute_end_at` bigint DEFAULT '0',
  `grp_display_name` varchar(100) DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_appkey_grpid_memid` (`app_key`,`group_id`,`member_id`),
  KEY `idx_memberid` (`app_key`,`member_id`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_type` tinyint DEFAULT '0',
  `user_id` varchar(32) NOT NULL,
  `nickname` varchar(50) DEFAULT NULL,
  `user_portrait` varchar(200) DEFAULT NULL,
  `pinyin` varchar(50) DEFAULT NULL,
  `phone` varchar(50) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `login_account` varchar(50) DEFAULT NULL,
  `login_pass` varchar(50) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_userid` (`app_key`,`user_id`),
  UNIQUE KEY `uniq_phone` (`app_key`,`phone`),
  UNIQUE KEY `uniq_email` (`app_key`,`email`),
  UNIQUE KEY `uniq_account` (`app_key`,`login_account`),
  KEY `idx_userid` (`app_key`,`user_type`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `userexts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `item_key` varchar(50) DEFAULT NULL,
  `item_value` varchar(2000) DEFAULT NULL,
  `item_type` tinyint DEFAULT '0',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_item_key` (`app_key`,`user_id`,`item_key`),
  KEY `idx_item_key` (`app_key`,`item_key`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `friendapplications` (
  `id` int NOT NULL AUTO_INCREMENT,
  `recipient_id` varchar(32) DEFAULT NULL,
  `sponsor_id` varchar(32) DEFAULT NULL,
  `apply_time` bigint DEFAULT NULL,
  `status` tinyint DEFAULT '0',
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_apply` (`app_key`,`recipient_id`,`sponsor_id`),
  KEY `idx_recipient` (`app_key`,`recipient_id`,`apply_time`),
  KEY `idx_sponsor` (`app_key`,`sponsor_id`,`apply_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `grpapplications` (
  `id` int NOT NULL AUTO_INCREMENT,
  `group_id` varchar(32) DEFAULT NULL,
  `apply_type` tinyint DEFAULT '0',
  `sponsor_id` varchar(32) DEFAULT NULL,
  `recipient_id` varchar(32) DEFAULT NULL,
  `inviter_id` varchar(32) DEFAULT NULL,
  `operator_id` varchar(32) DEFAULT NULL,
  `apply_time` bigint DEFAULT '0',
  `status` tinyint DEFAULT '0',
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_apply` (`app_key`,`group_id`,`apply_type`,`sponsor_id`,`recipient_id`),
  KEY `idx_sponsor` (`app_key`,`apply_type`,`sponsor_id`,`apply_time`),
  KEY `idx_group` (`app_key`,`apply_type`,`group_id`,`apply_time`),
  KEY `idx_recipient` (`app_key`,`apply_type`,`recipient_id`,`apply_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `qrcoderecords` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code_id` varchar(50) DEFAULT NULL,
  `status` tinyint DEFAULT NULL,
  `created_time` bigint DEFAULT NULL,
  `user_id` varchar(32) DEFAULT NULL,
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_id` (`app_key`,`code_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `smsrecords` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `phone` varchar(50) DEFAULT NULL,
  `code` varchar(10) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_phone` (`app_key`,`phone`,`created_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `ai_engines` (
  `id` int NOT NULL AUTO_INCREMENT,
  `engine_type` tinyint DEFAULT '0',
  `engine_conf` varchar(5000) DEFAULT NULL,
  `status` tinyint DEFAULT '0',
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_appkey` (`app_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `assistant_prompts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT NULL,
  `prompts` varchar(2000) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_id` (`app_key`,`id`),
  KEY `idx_user` (`app_key`,`user_id`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `botconfs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `bot_id` varchar(32) NULL,
  `nickname` varchar(50) DEFAULT NULL,
  `bot_portrait` varchar(200) DEFAULT NULL,
  `description` varchar(500) DEFAULT NULL,
  `bot_type` tinyint DEFAULT '0',
  `bot_conf` varchar(2000) NULL,
  `status` tinyint DEFAULT '0',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uniq_botid` (`app_key`, `bot_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `telebots` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(50) DEFAULT NULL,
  `bot_name` varchar(50) DEFAULT NULL,
  `bot_token` varchar(200) DEFAULT NULL,
  `status` tinyint DEFAULT NULL,
  `app_key` varchar(20) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_bot` (`app_key`,`bot_token`),
  KEY `idx_user` (`app_key`,`user_id`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `posts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `post_id` varchar(32) DEFAULT NULL,
  `title` varchar(200) DEFAULT NULL,
  `content` mediumblob,
  `content_brief` varchar(5000) DEFAULT NULL,
  `is_delete` tinyint DEFAULT '0',
  `user_id` varchar(32) DEFAULT NULL,
  `post_exset` mediumblob,
  `created_time` bigint DEFAULT '0',
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `status` tinyint DEFAULT '0',
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `uniq_id` (`app_key`,`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `postcomments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `comment_id` varchar(32) DEFAULT NULL,
  `post_id` varchar(32) DEFAULT NULL,
  `parent_comment_id` varchar(32) DEFAULT NULL,
  `parent_user_id` varchar(32) DEFAULT NULL,
  `user_id` varchar(32) DEFAULT NULL,
  `text` varchar(5000) DEFAULT NULL,
  `created_time` bigint DEFAULT NULL,
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `is_delete` tinyint DEFAULT '0',
  `status` tinyint DEFAULT '0',
  `app_key` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_id` (`app_key`,`comment_id`),
  KEY `idx_post` (`app_key`,`post_id`,`created_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT IGNORE INTO `globalconfs` (`conf_key`,`conf_value`)VALUES('jchatdb_versaion','20240716');