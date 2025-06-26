package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gmcc94/attendance-go/db"
	"github.com/go-chi/chi/v5"
)

func CreateAttendanceHandler(attendanceStore db.AttendanceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		studentIDStr := chi.URLParam(r, "id")
		studentID, err := strconv.Atoi(studentIDStr)
		if err != nil {
			http.Error(w, "Invalid student ID", http.StatusBadRequest)
			log.Printf("Error with student ID: %v", err)
			return
		}

		loc, err := time.LoadLocation("Europe/Dublin")
		if err != nil {
			log.Printf("Failed to load timezone: %v", err)
			http.Error(w, "Server timezone error", http.StatusInternalServerError)
			return
		}

		now := time.Now().In(loc)
		dateOnly := now.Truncate(24 * time.Hour)
		currentDay := now.Weekday().String()

		allowedDays := map[string]bool{
			"Monday":    true,
			"Tuesday":   true,
			"Wednesday": true,
			"Friday":    true,
		}

		if !allowedDays[currentDay] {
			http.Error(w, "No class on this day", http.StatusForbidden)
			return
		}

		// Insert attendance
		err = attendanceStore.InsertAttendance(studentID, dateOnly, currentDay)
		if err != nil {
			log.Printf("Failed to add attendance: %v", err)
			http.Error(w, "Failed to add attendance", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Attendance recorded successfully",
			"date":    now.Format("02/01/2006"),
			"day":     currentDay,
		})
	}
}

func GetStudentAttendanceByIDHandler(attendanceStore db.AttendanceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		studentIDStr := chi.URLParam(r, "id")
		studentID, err := strconv.Atoi(studentIDStr)
		if err != nil {
			http.Error(w, "Invalid student ID", http.StatusBadRequest)
			log.Printf("error with student ID %v", err)
			return
		}

		attendanceResp, err := attendanceStore.GetStudentAttendanceByID(studentID)
		if err != nil {
			http.Error(w, "Failed to fetch student attendance", http.StatusInternalServerError)
			log.Printf("Failed to fetch attendance: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(attendanceResp)
	}
}
