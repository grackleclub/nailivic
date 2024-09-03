CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    token TEXT NOT NULL,
    created_on TIMESTAMP,
    expires_on TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '1 hour'
);