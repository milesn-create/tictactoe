CREATE TABLE IF NOT EXISTS refresh_tokens(
    jti UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    used_at TIMESTAMPTZ
);