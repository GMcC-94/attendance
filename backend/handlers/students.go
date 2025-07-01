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

func UpdateStudentHandler(studentStore db.StudentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		studentID, err := helpers.GetStudentURLID(r)
		if err != nil {
			http.Error(w, "Invalid student ID", http.StatusBadRequest)
			return
		}

		var req struct {
			Name      *string `json:"name"`
			BeltGrade *string `json:"beltGrade"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			log.Printf("invalid JSON body %v", err)
			return
		}

		if req.Name == nil && req.BeltGrade == nil {
			http.Error(w, "No fields provided to update", http.StatusBadRequest)
			return
		}

		_, err = studentStore.UpdateStudent(studentID, req.Name, req.BeltGrade)
		if err != nil {
			http.Error(w, "Failed to update student: "+err.Error(), http.StatusInternalServerError)
			return
		}

		updatedStudent, err := studentStore.GetStudentByID(studentID)
		if err != nil {
			http.Error(w, "Failed to fetch updated student: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedStudent)
	}
}

func DeleteStudentHandler(studentStore db.StudentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentID, err := helpers.GetStudentURLID(r)
		if err != nil {
			log.Printf("invalid student ID %w", err)
			http.Error(w, "Invalid student ID", http.StatusBadRequest)
			return
		}

		err = studentStore.DeleteStudent(studentID)
		if err != nil {
			http.Error(w, "Failed to delete student: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
