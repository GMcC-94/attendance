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
	_, err := p.DB.Exec("INSERT into attendanceS (student_id, attendance_date, class_day) VALUES ($1, $2, $3)",
		studentID,
		date,
		dayOfWeek)
	return err
}
