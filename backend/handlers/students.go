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

		helpers.DecodeJSON(r, w, &req)

		dob, err := time.Parse("02/01/2006", req.DateOfBirth)
		if err != nil {
			http.Error(w, "Invalid date format, use DD/MM/YYYY", http.StatusBadRequest)
			return
		}

		err = studentStore.CreateStudent(req.Name, req.BeltGrade, req.StudentType, dob)
		if err != nil {
			log.Printf("Failed to create student %s", err)
			http.Error(w, "Failed to create student", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"message": "Student created successfully",
		}

		helpers.WriteJSON(w, http.StatusOK, response)
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

		helpers.WriteJSON(w, http.StatusOK, response)
	}
}

func GetAllAdultStudentsHandler(studentStore db.StudentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := studentStore.GetAllAdultStudents()
		if err != nil {
			log.Printf("failed to fetch students that are adults: %v", err)
			http.Error(w, "Failed to fetch all students", http.StatusInternalServerError)
			return
		}

		var response []types.StudentResponse
		for _, s := range students {
			response = append(response, types.StudentResponse{
				ID:          s.ID,
				Name:        s.Name,
				BeltGrade:   s.BeltGrade,
				DOB:         s.DateOfBirth.Format("02/01/2006"),
				Age:         helpers.CalculateAge(s.DateOfBirth),
				StudentType: s.StudentType,
			})
		}

		helpers.WriteJSON(w, http.StatusOK, response)
	}
}

func GetAllKidStudentsHandler(studentStore db.StudentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := studentStore.GetAllKidStudents()
		if err != nil {
			log.Printf("failed to fetch students that are kids: %v", err)
			http.Error(w, "Failed to fetch all students", http.StatusInternalServerError)
			return
		}

		var response []types.StudentResponse
		for _, s := range students {
			response = append(response, types.StudentResponse{
				ID:          s.ID,
				Name:        s.Name,
				BeltGrade:   s.BeltGrade,
				DOB:         s.DateOfBirth.Format("02/01/2006"),
				Age:         helpers.CalculateAge(s.DateOfBirth),
				StudentType: s.StudentType,
			})
		}

		helpers.WriteJSON(w, http.StatusOK, response)
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

		helpers.WriteJSON(w, http.StatusOK, updatedStudent)
	}
}

func DeleteStudentHandler(studentStore db.StudentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentID, err := helpers.GetStudentURLID(r)
		if err != nil {
			log.Printf("invalid student ID %v", err)
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
