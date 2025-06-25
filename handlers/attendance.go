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
			return
		}

		var req struct {
			AttendedDays []string `json:"attendedDays"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		allowedDays := map[string]bool{
			"Monday":    true,
			"Tuesday":   true,
			"Wednesday": true,
			"Friday":    true,
		}

		// TODO: move to a helper func as code will be reused
		loc, err := time.LoadLocation("Europe/Dublin")
		if err != nil {
			log.Printf("Failed to load timezone: %v", err)
			http.Error(w, "Server timezone error", http.StatusInternalServerError)
			return
		}

		now := time.Now().In(loc)
		currentDay := now.Weekday().String()

		for _, day := range req.AttendedDays {
			if !allowedDays[day] {
				http.Error(w, "Attendance cannot be taken for the current day", http.StatusForbidden)
				return
			}

			err = attendanceStore.InsertAttendance(studentID, now, currentDay)
			if err != nil {
				log.Printf("Failed to add attendance: %v", err)
				http.Error(w, "Failed to add attendance", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Attendance taken successfully",
			"date":    now.Format("02/01/2006"),
		})
	}
}
