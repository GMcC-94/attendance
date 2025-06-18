package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gmcc94/attendance-go/db"
)

func AddAttendanceHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Extract student ID from URL
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		studentIDStr := pathParts[2]
		studentID, err := strconv.Atoi(studentIDStr)
		if err != nil {
			http.Error(w, "Invalid student ID", http.StatusBadRequest)
			return
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

		allowedDays := map[string]bool{
			"Monday":    true,
			"Tuesday":   true,
			"Wednesday": true,
			"Friday":    true,
		}

		if !allowedDays[currentDay] {
			http.Error(w, "Attendance cannot be taken for the current day", http.StatusForbidden)
			return
		}

		err = db.InsertAttendance(database, studentID, now, currentDay)
		if err != nil {
			log.Printf("Failed to add attendance: %v", err)
			http.Error(w, "Failed to add attendance", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Attendance taken successfully",
			"date":    now.Format("02/01/2006"),
			"day":     currentDay,
		})
	}
}
