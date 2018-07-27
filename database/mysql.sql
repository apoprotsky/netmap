CREATE TABLE map.nodes (
  id int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  type enum ('point', 'switch') NOT NULL DEFAULT 'point',
  x mediumint(9) UNSIGNED DEFAULT 100,
  y mediumint(9) UNSIGNED DEFAULT 100,
  device smallint(6) UNSIGNED DEFAULT NULL,
  street smallint(6) UNSIGNED DEFAULT NULL,
  building varchar(8) DEFAULT NULL,
  porch varchar(4) DEFAULT NULL,
  level varchar(4) DEFAULT NULL,
  ip int(11) UNSIGNED DEFAULT NULL,
  base_node int(11) UNSIGNED DEFAULT NULL,
  base_port smallint(6) UNSIGNED DEFAULT NULL,
  uplink_port smallint(6) UNSIGNED DEFAULT NULL,
  notifications tinyint(1) UNSIGNED NOT NULL DEFAULT 0,
  status smallint(6) UNSIGNED NOT NULL DEFAULT 0,
  changed int(11) UNSIGNED NOT NULL DEFAULT 0,
  listed int(11) UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  INDEX IDX_nodes_base_node (base_node),
  INDEX IDX_nodes_type (type)
)
ENGINE = INNODB
CHARACTER SET utf8
COLLATE utf8_general_ci;

CREATE TABLE map.links (
  id int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `from` int(11) UNSIGNED NOT NULL,
  `to` int(11) UNSIGNED NOT NULL,
  PRIMARY KEY (id),
  INDEX idx_links_from (`from`),
  INDEX idx_links_to (`to`),
  CONSTRAINT fk_links_from_nodes_id FOREIGN KEY (`from`)
  REFERENCES map.nodes (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_links_to_nodes_id FOREIGN KEY (`to`)
  REFERENCES map.nodes (id) ON DELETE CASCADE ON UPDATE CASCADE
)
ENGINE = INNODB
CHARACTER SET utf8
COLLATE utf8_general_ci;
