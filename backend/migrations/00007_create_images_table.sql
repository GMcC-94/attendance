-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    type TEXT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
