package db

import (
	"database/sql"
	"time"
)

type AttendanceStore interface {
	InsertAttendance(studentID int, date time.Time, dayOfWeek string) error
}

type PostgresAttendanceStore struct {
	DB *sql.DB
}

func (p *PostgresAttendanceStore) InsertAttendance(studentID int, date time.Time, dayOfWeek string) error {
	_, err := p.DB.Exec("INSERT into attendance (student_id, date, day_of_week) VALUES ($1, $2, $2)",
		studentID,
		date,
		dayOfWeek)
	return err
}
