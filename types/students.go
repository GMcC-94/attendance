package types

import "time"

type Students struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	BeltGrade   string    `json:"beltGrade"`
	DateOfBirth time.Time `json:"dob"`
}

type CreateStudentRequest struct {
	Name        string `json:"name"`
	BeltGrade   string `json:"beltGrade"`
	DateOfBirth string `json:"dob"`
}

type StudentResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BeltGrade string `json:"beltGrade"`
	DOB       string `json:"dob"`
	Age       int    `json:"age"`
}

type StudentAttendance struct {
}
