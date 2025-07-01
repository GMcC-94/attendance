-- +goose Up
-- +goose StatementBegin
CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    context TEXT NOT NULL,                
    context_id INT,                        
    filename TEXT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd