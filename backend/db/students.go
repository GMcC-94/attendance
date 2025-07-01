package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gmcc94/attendance-go/types"
)

type StudentStore interface {
	CreateStudent(name, beltGrade, studentType string, dateofBirth time.Time) error
	GetAllStudents() ([]types.Students, error)
	GetAllAdultStudents() ([]types.Students, error)
	GetAllKidStudents() ([]types.Students, error)
	GetStudentByID(studentID int) (types.Students, error)
	UpdateStudent(studentID int, name, beltGrade *string) (types.Students, error)
	DeleteStudent(studentID int) error
}

type PostgresStudentStore struct {
	DB *sql.DB
}

func (p *PostgresStudentStore) CreateStudent(name, beltGrade, studentType string, dateOfBirth time.Time) error {
	_, err := p.DB.Exec(`
	INSERT INTO students (name, belt_grade, dob, student_type) 
	VALUES ($1, $2, $3, $4)`,
		name,
		beltGrade,
		dateOfBirth,
		studentType)

	return err
}

func (p *PostgresStudentStore) GetAllStudents() ([]types.Students, error) {
	rows, err := p.DB.Query("SELECT id, name, belt_grade, dob FROM students ORDER BY name ASC;")
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

func (p *PostgresStudentStore) GetAllAdultStudents() ([]types.Students, error) {
	rows, err := p.DB.Query(`
	SELECT id, name, belt_grade, student_type, dob 
	FROM students 
	WHERE student_type = 'adult'
	ORDER BY name ASC;
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Students
	for rows.Next() {
		var s types.Students
		if err := rows.Scan(&s.ID, &s.Name, &s.BeltGrade, &s.StudentType, &s.DateOfBirth); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

func (p *PostgresStudentStore) GetAllKidStudents() ([]types.Students, error) {
	rows, err := p.DB.Query(`
	SELECT id, name, belt_grade, student_type, dob 
	FROM students 
	WHERE student_type = 'kid'
	ORDER BY name ASC;
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Students
	for rows.Next() {
		var s types.Students
		if err := rows.Scan(&s.ID, &s.Name, &s.BeltGrade, &s.StudentType, &s.DateOfBirth); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

func (p *PostgresStudentStore) GetStudentByID(studentID int) (types.Students, error) {
	var student types.Students

	row := p.DB.QueryRow(`
	SELECT name, belt_grade, dob
	FROM students
	WHERE id = $1;`, studentID)
	err := row.Scan(&student.Name, &student.BeltGrade, &student.DateOfBirth)
	if errors.Is(err, sql.ErrNoRows) {
		return types.Students{}, fmt.Errorf("query student by id: %w", err)
	}

	return student, nil
}

func (p *PostgresStudentStore) UpdateStudent(studentID int, name, beltGrade *string) (types.Students, error) {
	setClauses := []string{}
	args := []interface{}{studentID}
	argPos := 2

	if name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *name)
		argPos++
	}
	if beltGrade != nil {
		setClauses = append(setClauses, fmt.Sprintf("belt_grade = $%d", argPos))
		args = append(args, *beltGrade)
		argPos++
	}

	if len(setClauses) == 0 {
		return types.Students{}, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE students SET %s WHERE id = $1", strings.Join(setClauses, ", "))

	_, err := p.DB.Exec(query, args...)
	if err != nil {
		log.Printf("error updating student: %v", err)
		return types.Students{}, err
	}

	return types.Students{ID: studentID}, nil
}

func (p *PostgresStudentStore) DeleteStudent(studentID int) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
	DELETE FROM attendances
	WHERE student_id = $1`, studentID)
	if err != nil {
		return fmt.Errorf("failed to delete attendances: %w", err)
	}

	res, err := tx.Exec(`
	DELETE FROM students
	WHERE id = $1`, studentID)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not determine affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no student found with ID %d", studentID)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
