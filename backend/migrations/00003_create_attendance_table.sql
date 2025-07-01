-- +goose Up
CREATE TABLE IF NOT EXISTS attendances (
    id SERIAL PRIMARY KEY,
    student_id INT NOT NULL REFERENCES students(id),
    class_day VARCHAR(20) NOT NULL,
    attendance_date TIMESTAMP NOT NULL
);


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS attendances;
-- +goose StatementEnd
