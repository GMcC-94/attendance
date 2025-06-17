package db_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	studentDB "github.com/gmcc94/attendance-go/db"
)

func TestCreateStudent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock DB: %s", err)
	}
	defer db.Close()

	name := "Gerard McCann"
	beltGrade := "Green Belt"
	dob := time.Date(1994, 7, 19, 0, 0, 0, 0, time.Local)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO students (name, belt_grade, dob) VALUES ($1, $2, $3)")).
		WithArgs(name, beltGrade, dob).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = studentDB.CreateStudent(db, name, beltGrade, dob)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectation: %v", err)
	}

}
