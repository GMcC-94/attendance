-- +goose Up
-- +goose StatementBegin
CREATE TABLE club_accounts (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    amount NUMERIC(12,2) NOT NULL CHECK (amount >= 0),
    category TEXT NOT NULL CHECK (category IN ('income', 'expenditure')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_club_accounts_created_at ON club_accounts(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS club_accounts;
-- +goose StatementEnd
