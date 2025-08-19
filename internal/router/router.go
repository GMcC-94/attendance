package router

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gmcc94/attendance-go/internal/config"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// New creates and configures the application router
func New(db *sql.DB, cfg *config.Config) chi.Router {
	r := chi.NewRouter()

	// Global middleware
	setupMiddleware(r)

	// Routes
	setupRoutes(r, db, cfg)

	return r
}

func setupMiddleware(r chi.Router) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(30 * time.Second))
}

func setupRoutes(r chi.Router, db *sql.DB, cfg *config.Config) {
	// Health check
	r.Get("/health", healthCheckHandler)

	// Temporary root route
	r.Get("/", homeHandler)

	// TODO: Add more routes as we build out the application
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
		<h1>🥋 Attendance System</h1>
		<p>Server is running successfully!</p>
		<p>Ready for development...</p>
		<ul>
			<li><a href="/health">Health Check</a></li>
		</ul>
	`))
}
