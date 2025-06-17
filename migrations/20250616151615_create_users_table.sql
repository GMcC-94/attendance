-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS users{
    id SERIAL PRIMARY KEY,
    username varchar(100) NOT NULL,
    password_hash TEXT NOT NULL
}

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS users;
