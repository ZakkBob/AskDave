DROP TABLE IF EXISTS crawl;
DROP TABLE IF EXISTS link;
DROP TABLE IF EXISTS page;
DROP TABLE IF EXISTS robots;
DROP TABLE IF EXISTS site;
DROP TYPE IF EXISTS failure_reason;

CREATE TABLE site (
  id SERIAL PRIMARY KEY NOT NULL,
  url varchar(200) NOT NULL UNIQUE
);

CREATE TABLE page (
  id SERIAL PRIMARY KEY NOT NULL,

  site integer references site(id) NOT NULL,
  path varchar(2048) NOT NULL,
  
  title varchar(50),
  og_title varchar(50),
  og_description varchar(50),
  og_sitename varchar(50),

  hash char(32),

  next_crawl date,
  crawl_interval integer,
  interval_delta integer,

  assigned bool DEFAULT FALSE NOT NULL,

  UNIQUE (site, path)
);

CREATE TYPE failure_reason AS ENUM ('NoFailure', 'RobotsDisallowed', 'FetchFailed');

CREATE TABLE crawl (
  id SERIAL PRIMARY KEY NOT NULL,
  page integer references page(id) NOT NULL,
  datetime timestamp NOT NULL,
  success bool NOT NULL,
  failure_reason failure_reason,
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
  site integer references site(id) NOT NULL,
  allowed_patterns varchar(50)[],
  disallowed_patterns varchar(50)[],
  last_crawl date
);