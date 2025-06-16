package main

import (
	"log"
	"net/http"

	"github.com/gmcc94/attendance-go/config"
	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/handlers"
	"github.com/gorilla/mux"
)

func main() {

	config.LoadConfig()

	db := db.InitDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/signup", handlers.SignupHandler(db)).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler(db)).Methods("POST")

	log.Println("Server starting on port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
