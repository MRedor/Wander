CREATE TABLE IF NOT EXISTS  `lists` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `type` enum('routes','objects') NOT NULL,
    `name` text,
    `views` int(11) NOT NULL DEFAULT 1,
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `objects_in_list` (
    `list_id` int(11) unsigned NOT NULL,
    `object_id` int(11) unsigned NOT NULL,
    PRIMARY KEY (list_id, object_id),
    FOREIGN KEY (list_id) REFERENCES lists(id),
    FOREIGN KEY (object_id) REFERENCES points(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `routes_in_list` (
    `list_id` int(11) unsigned NOT NULL,
    `route_id` int(11) unsigned NOT NULL,
    PRIMARY KEY (list_id, route_id),
    FOREIGN KEY (list_id) REFERENCES lists(id),
    FOREIGN KEY (route_id) REFERENCES routes(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;