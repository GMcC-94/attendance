package types

import "time"

type Students struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	BeltGrade   string    `json:"beltGrade"`
	DateOfBirth time.Time `json:"dob"`
	StudentType string    `json:"studentType"`
}

type CreateStudentRequest struct {
	Name        string `json:"name"`
	BeltGrade   string `json:"beltGrade"`
	DateOfBirth string `json:"dob"`
	StudentType string `json:"studentType"`
}

type StudentResponse struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	BeltGrade   string              `json:"beltGrade"`
	DOB         string              `json:"dob"`
	Age         int                 `json:"age"`
	StudentType string              `json"studentType"`
	Attendance  []StudentAttendance `json:"attendance"`
}

type StudentAttendance struct {
	Date     string `json:"date"`
	ClassDay string `json:"classDay"`
}

type StudentAttendanceResponse struct {
	ID         int                 `json:"id"`
	Name       string              `json:"name"`
	BeltGrade  string              `json:"beltGrade"`
	DOB        string              `json:"dob"`
	Age        int                 `json:"age"`
	Attendance []StudentAttendance `json:"attendance"`
}
