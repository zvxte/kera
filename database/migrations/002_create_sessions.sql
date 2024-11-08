CREATE TABLE IF NOT EXISTS sessions(
    id BYTEA NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    creation_date DATE NOT NULL,
    expiration_date DATE NOT NULL
);
