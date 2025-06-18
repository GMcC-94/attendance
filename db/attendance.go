package db

import (
	"database/sql"
	"time"
)

func InsertAttendance(db *sql.DB, studentID int, date time.Time, dayOfWeek string) error {
	_, err := db.Exec("INSERT into attendance (student_id, date, day_of_week) VALUES ($1, $2, $2)",
		studentID,
		date,
		dayOfWeek)
	return err
}
