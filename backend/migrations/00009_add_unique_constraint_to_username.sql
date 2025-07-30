-- +goose Up
-- +goose StatementBegin
-- Remove duplicates (keep the lowest id per username)
DELETE FROM users a
USING users b
WHERE a.username = b.username
  AND a.id > b.id;

-- Add unique constraint
ALTER TABLE users
ADD CONSTRAINT users_username_unique UNIQUE (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the unique constraint
ALTER TABLE users
DROP CONSTRAINT IF EXISTS users_username_unique;
-- +goose StatementEnd
