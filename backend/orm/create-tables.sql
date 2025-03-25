DROP TABLE IF EXISTS crawl;
DROP TABLE IF EXISTS link;
DROP TABLE IF EXISTS page;
DROP TABLE IF EXISTS robots;
DROP TABLE IF EXISTS site;
DROP TYPE IF EXISTS failure_reason;

CREATE TABLE site (
  id SERIAL PRIMARY KEY NOT NULL,
  url varchar(200) NOT NULL UNIQUE,

  allowed_patterns varchar(100)[],
  disallowed_patterns varchar(100)[],
  last_robots_crawl timestamp
);

CREATE TABLE page (
  id SERIAL PRIMARY KEY NOT NULL,

  site integer references site(id) NOT NULL,
  path varchar(2048) NOT NULL,
  
  title varchar(50) DEFAULT '' NOT NULL,
  og_title varchar(50) DEFAULT '' NOT NULL,
  og_description varchar(100) DEFAULT '' NOT NULL,
  og_sitename varchar(50) DEFAULT '' NOT NULL,

  hash char(32) NOT NULL,

  next_crawl date NOT NULL,
  crawl_interval integer DEFAULT 7 CHECK(crawl_interval > 0),
  interval_delta integer DEFAULT 1,

  assigned bool DEFAULT FALSE NOT NULL,

  UNIQUE (site, path)
);

CREATE TYPE failure_reason AS ENUM ('NoFailure', 'RobotsDisallowed', 'FetchFailed');

-- CREATE TABLE crawl (
--   id SERIAL PRIMARY KEY NOT NULL,
--   page integer references page(id) NOT NULL,
  
--   datetime timestamp NOT NULL,
--   success bool NOT NULL,
--   failure_reason int NOT NULL,
--   content_changed bool NOT NULL,
--   hash char(32) NOT NULL
-- );

CREATE TABLE link (
  id SERIAL PRIMARY KEY NOT NULL,
  src integer references page(id) NOT NULL,
  dst integer references page(id) NOT NULL
);