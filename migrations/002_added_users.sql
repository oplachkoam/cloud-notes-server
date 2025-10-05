CREATE TYPE USER_STATUS AS ENUM ('pending', 'active', 'blocked', 'deleted');

CREATE TABLE IF NOT EXISTS users
(
    id            UUID PRIMARY KEY,
    login         TEXT UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    first_name    TEXT        NOT NULL,
    timezone      TEXT        NOT NULL,
    status        USER_STATUS NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL
);
