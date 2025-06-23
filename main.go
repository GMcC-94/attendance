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

	db := db.InitDB()
	defer db.Close()

	r := chi.NewRouter()

	// Auth Routes
	r.Post("/signup", handlers.SignupHandler(db))
	r.Post("/login", handlers.LoginHandler(db))

	// Student Routes
	r.Post("/students", handlers.CreateStudentHandler(db))
	r.Get("/students", handlers.GetAllStudentsHandler(db))

	log.Println("Server starting on port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
