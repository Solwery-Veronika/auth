CREATE TABLE IF NOT EXISTS participants (
    id BIGSERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

DROP TABLE IF EXISTS participants;