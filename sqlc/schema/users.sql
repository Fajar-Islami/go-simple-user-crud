CREATE TABLE IF NOT EXISTS users (
  id   BIGINT  NOT NULL PRIMARY KEY,
  email varchar(255) unique not null,
  first_name varchar(255) collate utf8_general_ci unique not null,
  last_name varchar(255) collate utf8_general_ci unique not null,
  avatar text,
  created_at timestamp default now(),
  updated_at timestamp,
  deleted_at timestamp,
  UNIQUE KEY unique_name (first_name collate utf8_general_ci, last_name collate utf8_general_ci)
);
