-- +goose Up
-- +goose StatementBegin
ALTER TABLE students
ADD COLUMN student_type VARCHAR(20);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE students
DROP COLUMN student_type;
-- +goose StatementEnd
