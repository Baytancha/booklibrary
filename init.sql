-- gobook/init.sql
\c gobook;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name  text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    deleted bool NOT NULL DEFAULT false,
    version integer NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS authors (
    id bigserial PRIMARY KEY,
    name  VARCHAR(255) NOT NULL,
    times_ordered integer NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS books (
    id bigserial PRIMARY KEY,
    year integer NOT NULL,
    title VARCHAR(255) NOT NULL,
    available bool NOT NULL DEFAULT true,
    author_id SERIAL REFERENCES authors(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS rented (
  id bigserial PRIMARY KEY,
  book_id bigserial UNIQUE REFERENCES books(id),
  user_id bigserial REFERENCES users(id),
  active bool NOT NULL DEFAULT true
);
