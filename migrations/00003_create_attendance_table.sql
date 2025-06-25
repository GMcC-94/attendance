-- +goose Up
CREATE TABLE IF NOT EXISTS attendances (
    id SERIAL PRIMARY KEY,
    student_id INT NOT NULL REFERENCES students(id),
    class_day VARCHAR(20) NOT NULL,
    attendance_date DATE NOT NULL
);


-- +goose Down
DROP TABLE IF EXISTS attendance;
