package db_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/types"
)

func TestCreateStudent(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock DB: %s", err)
	}
	defer sqlDB.Close()

	studentStore := &db.PostgresStudentStore{DB: sqlDB}

	name := "Gerard McCann"
	beltGrade := "Green Belt"
	dob := time.Date(1994, 7, 19, 0, 0, 0, 0, time.Local)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO students (name, belt_grade, dob) VALUES ($1, $2, $3)")).
		WithArgs(name, beltGrade, dob).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = studentStore.CreateStudent(name, beltGrade, dob)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectation: %v", err)
	}
}

func TestGetAllStudents(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock DB: %s", err)
	}
	defer sqlDB.Close()

	studentStore := &db.PostgresStudentStore{DB: sqlDB}

	rows := sqlmock.NewRows([]string{"id", "name", "belt_grade", "dob"}).
		AddRow(1, "John McTasney", "6th Dan Black Belt", time.Date(1994, 7, 19, 0, 0, 0, 0, time.Local)).
		AddRow(2, "Pearse Rice", "Black Belt", time.Date(2004, 5, 1, 0, 0, 0, 0, time.Local))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, belt_grade, dob FROM students")).
		WillReturnRows(rows)

	students, err := studentStore.GetAllStudents()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if len(students) != 2 {
		t.Errorf("Expected 2 students, got %d", len(students))
	}

	expected := []types.Students{
		{ID: 1, Name: "John McTasney", BeltGrade: "6th Dan Black Belt", DateOfBirth: time.Date(1994, 7, 19, 0, 0, 0, 0, time.Local)},
		{ID: 2, Name: "Pearse Rice", BeltGrade: "Black Belt", DateOfBirth: time.Date(2004, 5, 1, 0, 0, 0, 0, time.Local)},
	}

	for i, student := range students {
		if student != expected[i] {
			t.Errorf("Expected student %v, got %v", expected[i], student)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}
}
