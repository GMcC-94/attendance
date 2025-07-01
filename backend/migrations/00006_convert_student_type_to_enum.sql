-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
  CREATE TYPE student_type_enum AS ENUM ('kid', 'adult');
EXCEPTION
  WHEN duplicate_object THEN null;
END $$;

ALTER TABLE students
  ALTER COLUMN student_type TYPE student_type_enum
  USING student_type::student_type_enum;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE students
  ALTER COLUMN student_type TYPE VARCHAR(20);

DROP TYPE IF EXISTS student_type_enum;
-- +goose StatementEnd
