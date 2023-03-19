CREATE TABLE IF NOT EXISTS users (
  id   BIGINT UNSIGNED  NOT NULL PRIMARY KEY,
  email varchar(255) unique not null,
  first_name varchar(255) collate utf8_general_ci not null,
  last_name varchar(255) collate utf8_general_ci not null,
  avatar text,
  created_at timestamp default now(),
  updated_at timestamp,
  deleted_at timestamp
);
