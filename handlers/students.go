package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/helpers"
	"github.com/gmcc94/attendance-go/types"
)

func CreateStudentHandler(studentStore db.StudentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateStudentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		dob, err := time.Parse("02/01/2006", req.DateOfBirth)
		if err != nil {
			http.Error(w, "Invalid date format, use DD/MM/YYYY", http.StatusBadRequest)
			return
		}

		err = studentStore.CreateStudent(req.Name, req.BeltGrade, dob)
		if err != nil {
			log.Printf("Failed to create student %s", err)
			http.Error(w, "Failed to create student", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"message": "Student created successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetAllStudentsHandler(studentStore db.StudentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := studentStore.GetAllStudents()
		if err != nil {
			http.Error(w, "Failed to fetch all students", http.StatusInternalServerError)
			return
		}

		var response []types.StudentResponse
		for _, s := range students {
			response = append(response, types.StudentResponse{
				ID:        s.ID,
				Name:      s.Name,
				BeltGrade: s.BeltGrade,
				DOB:       s.DateOfBirth.Format("02/01/2006"),
				Age:       helpers.CalculateAge(s.DateOfBirth),
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
