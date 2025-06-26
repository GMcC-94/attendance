package db

import (
	"database/sql"
	"time"

	"github.com/gmcc94/attendance-go/helpers"
	"github.com/gmcc94/attendance-go/types"
)

type AttendanceStore interface {
	InsertAttendance(studentID int, date time.Time, dayOfWeek string) error
	GetStudentAttendanceByID(studentID int) (types.StudentAttendanceResponse, error)
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

func (p *PostgresAttendanceStore) GetStudentAttendanceByID(studentID int) (types.StudentAttendanceResponse, error) {
	rows, err := p.DB.Query(`
	SELECT 
    s.id, s.name, s.belt_grade, s.dob,
    a.attendance_date, a.class_day
	FROM students s
	LEFT JOIN attendances a ON s.id = a.student_id
	WHERE s.id = $1
	ORDER BY a.attendance_date ASC;
	`, studentID)
	if err != nil {
		return types.StudentAttendanceResponse{}, err
	}
	defer rows.Close()

	var resp types.StudentAttendanceResponse
	var attendanceList []types.StudentAttendance

	for rows.Next() {
		var (
			id        int
			name      string
			beltGrade string
			dob       time.Time
			date      sql.NullTime
			classDay  sql.NullString
		)

		if err := rows.Scan(&id, &name, &beltGrade, &dob, &date, &classDay); err != nil {
			return types.StudentAttendanceResponse{}, err
		}

		resp.ID = id
		resp.Name = name
		resp.BeltGrade = beltGrade
		resp.DOB = dob.Format("02/01/2006")
		resp.Age = helpers.CalculateAge(dob)

		if date.Valid && classDay.Valid {
			attendanceList = append(attendanceList, types.StudentAttendance{
				Date:     date.Time.Format("02/01/2006"),
				ClassDay: classDay.String,
			})
		}
	}

	resp.Attendance = attendanceList
	return resp, nil
}
