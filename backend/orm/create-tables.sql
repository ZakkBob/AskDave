DROP TABLE IF EXISTS crawl;
DROP TABLE IF EXISTS link;
DROP TABLE IF EXISTS page;
DROP TABLE IF EXISTS robots;
DROP TABLE IF EXISTS site;

CREATE TABLE site (
  id SERIAL PRIMARY KEY NOT NULL,
  url varchar(200) NOT NULL UNIQUE
);

CREATE TABLE page (
  id SERIAL PRIMARY KEY NOT NULL,

  site integer references site(id) NOT NULL,
  path varchar(2048) NOT NULL,
  
  title varchar(50) DEFAULT '' NOT NULL,
  og_title varchar(50) DEFAULT '' NOT NULL,
  og_description varchar(500) DEFAULT '' NOT NULL,
  og_sitename varchar(50) DEFAULT '' NOT NULL,

  hash char(32) DEFAULT 'd41d8cd98f00b204e9800998ecf8427e' NOT NULL,

  next_crawl date,
  crawl_interval integer DEFAULT 7,
  interval_delta integer DEFAULT 1,

  assigned bool DEFAULT FALSE NOT NULL,

  UNIQUE (site, path)
);

CREATE TABLE crawl (
  id SERIAL PRIMARY KEY NOT NULL,
  page integer references page(id) NOT NULL,
  datetime timestamp NOT NULL,
  success bool NOT NULL,
  failure_reason int,
  content_changed bool,
  hash char(32)
);

CREATE TABLE link (
  id SERIAL PRIMARY KEY NOT NULL,
  src integer references page(id) NOT NULL,
  dst integer references page(id) NOT NULL
);

CREATE TABLE robots (
  id SERIAL PRIMARY KEY NOT NULL,
  site integer references site(id) NOT NULL UNIQUE,
  allowed_patterns varchar(50)[],
  disallowed_patterns varchar(50)[],
  last_crawl date
);