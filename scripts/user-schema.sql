\c user_db

CREATE TABLE users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    username      VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT        NOT NULL,
    expires_at TIMESTAMP   NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT NOW()
);