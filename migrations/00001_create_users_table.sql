-- +goose Up
CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    username varchar(100) NOT NULL,
    password_hash TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;
