package db

import (
	"database/sql"
	"time"

	"github.com/gmcc94/attendance-go/types"
)

func CreateStudent(db *sql.DB, name, beltGrade string, dateOfBirth time.Time) error {
	_, err := db.Exec("INSERT INTO users (name, belt_grade, dob) VALUES ($1, $2, $3)", name, beltGrade, dateOfBirth)

	return err
}

func GetAllStudents(db *sql.DB) ([]types.Students, error) {
	rows, err := db.Query("SELECT id, name, belt_grade, dob FROM students")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Students
	for rows.Next() {
		var s types.Students
		if err := rows.Scan(&s.ID, &s.Name, &s.BeltGrade, &s.DateOfBirth); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}
