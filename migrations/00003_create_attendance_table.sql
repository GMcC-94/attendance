-- +goose Up
CREATE TABLE IF NOT EXISTS attendance (
    id SERIAL PRIMARY KEY,
    student_id INT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    day_of_week VARCHAR(10) NOT NULL,
    attended BOOLEAN NOT NULL DEFAULT TRUE
);


-- +goose Down
DROP TABLE IF EXISTS attendance;
