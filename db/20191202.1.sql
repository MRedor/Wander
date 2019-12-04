CREATE TABLE IF NOT EXISTS  `routes` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `type` enum('direct','round') NOT NULL,
    `start_lat` float(10,7) NOT NULL,
    `start_lon` float(10,7) NOT NULL,
    `finish_lat` float(10,7) NOT NULL,
    `finish_lon` float(10,7) NOT NULL,
    `length` float(10,7) NOT NULL,
    `time` int(11) NOT NULL,
    `objects` text,
    `points` text,
    `name` text,
    `count` int(11) NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
