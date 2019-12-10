CREATE TABLE IF NOT EXISTS  `routes` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `type` enum('direct','round') NOT NULL,
    `start_lat` float(10,2) NOT NULL,
    `start_lon` float(10,2) NOT NULL,
    `finish_lat` float(10,2),
    `finish_lon` float(10,2),
    `radius` int(11),
    `length` float(10,7) NOT NULL,
    `time` int(11) NOT NULL,
    `objects` text,
    `points` text,
    `name` text,
    `count` int(11) NOT NULL DEFAULT 1,
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
