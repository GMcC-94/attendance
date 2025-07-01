package helpers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetStudentURLID(r *http.Request) (int, error) {
	studentIDStr := chi.URLParam(r, "id")
	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		return 0, fmt.Errorf("error extracting studentID from URL param %w", err)
	}
	return studentID, nil
}
