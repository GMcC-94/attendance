-- +goose Up
CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    belt_grade VARCHAR(50),
    dob DATE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS students;