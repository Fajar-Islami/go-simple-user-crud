CREATE TABLE IF NOT EXISTS users (
  id   BIGINT  NOT NULL PRIMARY KEY,
  email varchar(255) unique not null,
  first_name varchar(255) unique not null,
  last_name varchar(255) unique not null,
  avatar text,
  created_at timestamp default now(),
  updated_at timestamp,
  deleted_at timestamp
);
