CREATE TABLE `advertisement` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `start_at` datetime NOT NULL,
  `end_at` datetime NOT NULL
);

CREATE TABLE `gender` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(20) NOT NULL,
  `code` char(1) NOT NULL
);

CREATE TABLE `country` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `code` char(2) NOT NULL
);

CREATE TABLE `platform` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255) NOT NULL
);

CREATE TABLE `cond` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `age_start` int,
  `age_end` int
);

CREATE TABLE `advertisement_cond` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `advertisement_id` int NOT NULL,
  `cond_id` int NOT NULL
);

CREATE TABLE `cond_gender` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `cond_id` int NOT NULL,
  `gender_id` int NOT NULL
);

CREATE TABLE `cond_country` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `cond_id` int NOT NULL,
  `country_id` int NOT NULL
);

CREATE TABLE `cond_platform` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `cond_id` int NOT NULL,
  `platform_id` int NOT NULL
);

ALTER TABLE `advertisement_cond` ADD FOREIGN KEY (`advertisement_id`) REFERENCES `advertisement` (`id`);

ALTER TABLE `advertisement_cond` ADD FOREIGN KEY (`cond_id`) REFERENCES `cond` (`id`);

ALTER TABLE `cond_gender` ADD FOREIGN KEY (`cond_id`) REFERENCES `cond` (`id`);

ALTER TABLE `cond_gender` ADD FOREIGN KEY (`gender_id`) REFERENCES `gender` (`id`);

ALTER TABLE `cond_country` ADD FOREIGN KEY (`cond_id`) REFERENCES `cond` (`id`);

ALTER TABLE `cond_country` ADD FOREIGN KEY (`country_id`) REFERENCES `country` (`id`);

ALTER TABLE `cond_platform` ADD FOREIGN KEY (`cond_id`) REFERENCES `cond` (`id`);

ALTER TABLE `cond_platform` ADD FOREIGN KEY (`platform_id`) REFERENCES `platform` (`id`);