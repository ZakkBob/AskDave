DROP TABLE IF EXISTS album;
DROP TABLE IF EXISTS crawl;
DROP TABLE IF EXISTS link;
DROP TABLE IF EXISTS page;
DROP TABLE IF EXISTS site;
DROP TYPE IF EXISTS failure_reason;

CREATE TABLE site (
  id SERIAL PRIMARY KEY NOT NULL,
  url varchar(200) NOT NULL
);

CREATE TABLE page (
  id SERIAL PRIMARY KEY NOT NULL,
  site integer references site(id) NOT NULL,
  url varchar(2048) NOT NULL,
  next_crawl date,
  crawl_interval int,
  interval_delta int
);

CREATE TYPE failure_reason AS ENUM ('NoFailure', 'RobotsDisallowed', 'FetchFailed');

CREATE TABLE crawl (
  id SERIAL PRIMARY KEY NOT NULL,
  page integer references page(id) NOT NULL,
  datetime timestamp NOT NULL,
  success bool NOT NULL,
  failure_reason failure_reason,
  content_changed bool,
  title varchar(50),
  og_title varchar(50),
  og_description varchar(50),
  hash uuid
);

CREATE TABLE link (
  id SERIAL PRIMARY KEY NOT NULL,
  src integer references page(id) NOT NULL,
  dst varchar(2048) NOT NULL,
  count integer NOT NULL
)

-- INSERT INTO
--   album (title, artist, price)
-- VALUES
--   ('Blue Train', 'John Coltrane', 56.99),
--   ('Giant Steps', 'John Coltrane', 63.99),
--   ('Jeru', 'Gerry Mulligan', 17.99),
--   ('Sarah Vaughan', 'Sarah Vaughan', 34.98);