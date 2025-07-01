package main

import (
	"log"
	"net/http"

	"github.com/gmcc94/attendance-go/config"
	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {

	config.LoadConfig()

	sqlDB := db.InitDB()
	defer sqlDB.Close()

	userStore := &db.PostgresUserStore{DB: sqlDB}
	refreshTokenStore := &db.PostgresRefreshTokenStore{DB: sqlDB}
	studentStore := &db.PostgresStudentStore{DB: sqlDB}
	attendanceStore := &db.PostgresAttendanceStore{DB: sqlDB}
	imageStore := &db.PostgresImageStore{DB: sqlDB}

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		// Auth Routes
		r.Post("/signup", handlers.SignupHandler(userStore))
		r.Post("/login", handlers.LoginHandler(userStore, refreshTokenStore))

		// Student Routes
		r.Post("/students", handlers.CreateStudentHandler(studentStore))
		r.Get("/students", handlers.GetAllStudentsHandler(studentStore))
		r.Get("/students/adults", handlers.GetAllAdultStudentsHandler(studentStore))
		r.Get("/students/kids", handlers.GetAllKidStudentsHandler(studentStore))
		r.Put("/students/{id}", handlers.UpdateStudentHandler(studentStore))
		r.Delete("/students/{id}", handlers.DeleteStudentHandler(studentStore))

		// Attendance Routes
		r.Post("/students/{id}/attendance", handlers.CreateAttendanceHandler(attendanceStore))
		r.Get("/students/{id}/attendance", handlers.GetStudentAttendanceByIDHandler(attendanceStore))

		// Image upload
		r.Post("/upload", handlers.UploadImageHandler(imageStore))
	})

	fs := http.FileServer(http.Dir("./uploads"))
	r.Handle("/uploads/", http.StripPrefix("/uploads", fs))

	log.Println("Server starting on port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
