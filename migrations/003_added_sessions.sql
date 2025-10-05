CREATE TABLE IF NOT EXISTS sessions
(
    id         UUID PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users (id),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL
);
