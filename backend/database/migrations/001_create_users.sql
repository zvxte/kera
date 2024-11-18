CREATE TABLE IF NOT EXISTS users(
    id UUID NOT NULL PRIMARY KEY,
    username VARCHAR(16) NOT NULL,
    username_lower VARCHAR(16) UNIQUE NOT NULL,
    display_name VARCHAR(16) NOT NULL,
    hashed_password TEXT NOT NULL,
    creation_date DATE NOT NULL
);
