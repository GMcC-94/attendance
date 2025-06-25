package db

import (
	"database/sql"
	"time"

	"github.com/gmcc94/attendance-go/types"
)

type StudentStore interface {
	CreateStudent(name, beltGrade string, dateofBirth time.Time) error
	GetAllStudents() ([]types.Students, error)
}

type PostgresStudentStore struct {
	DB *sql.DB
}

func (p *PostgresStudentStore) CreateStudent(name, beltGrade string, dateOfBirth time.Time) error {
	_, err := p.DB.Exec("INSERT INTO students (name, belt_grade, dob) VALUES ($1, $2, $3)", name, beltGrade, dateOfBirth)

	return err
}

func (p *PostgresStudentStore) GetAllStudents() ([]types.Students, error) {
	rows, err := p.DB.Query("SELECT id, name, belt_grade, dob FROM students")
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
