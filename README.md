# Attendance Tracking Application

A simple and secure attendance tracking system built with Go, PostgreSQL, and Docker, designed for multi-tenant use and future AWS deployment.

---

## Features

- ✅ User authentication with JWT tokens  
- ✅ Student management including date of birth and belt grade  
- ✅ Age calculation based on UK timezone and DD/MM/YYYY date format  
- ✅ Database migrations using [Goose](https://github.com/pressly/goose)  
- ✅ Dockerized PostgreSQL database for easy setup  
- ✅ RESTful API endpoints for attendance operations  

> **Note:**  
> The current application is intended for internal use by **admin users only.**  
> User role differentiation (e.g., students, instructors) is not planned at this stage but may be considered in future versions.

---

## Planned Features

- [ ] Attendance reporting and summary dashboards
- [ ] Frontend interface for admin users
- [ ] Multi-tenant architecture to support multiple organizations or schools
- [ ] Deployment to AWS using services such as ECS, RDS, and S3
- [ ] CI/CD integration for automated testing and deployment
- [ ] API documentation with Swagger

---

## Tech Stack & Tools

| Tool            | Purpose                                  |
|-----------------|------------------------------------------|
| **Go**          | Backend API development                   |
| **PostgreSQL**  | Relational database to store users, students, and attendance data |
| **Docker**      | Containerization of PostgreSQL database  |
| **Goose**       | Database migration management             |
| **JWT**         | Secure user authentication                 |
| **sqlmock**     | Mocking database in unit tests             |

---

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (v1.20+ recommended)  
- [Docker](https://docs.docker.com/get-docker/)  
- [Goose CLI](https://github.com/pressly/goose) installed  
  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/attendance-tracker.git
   cd attendance-tracker
2. **Create .env file**
   ```bash
   Add environment variables (example):
   DB_URL=postgres://user:password@localhost:5432/attendanceDB?sslmode=disable
   JWT_SECRET=your-secure-jwt-secret-key
3. **Run PostgreSQL with Docker**
   ```bash
   Run make docker-up from CLI to start a docker container
4. **Run database migrations**
   ```bash
   Run database migrations with make migrate-up from CLI
