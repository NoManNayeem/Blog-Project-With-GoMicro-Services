package main

import (
	"database/sql"
	"log"
	"net/http"
	_ "user-management/docs" // For Swagger documentation
	"user-management/user"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	httpSwagger "github.com/swaggo/http-swagger"
)

var db *sql.DB

// @title User Management API
// @version 1.0
// @description This is a simple User Management API for handling user registration, login, and RBAC.
// @host localhost:8000
// @BasePath /

func main() {
	// MySQL connection string
	dsn := "root:@tcp(127.0.0.1:3306)/user_management" // Replace with your MySQL credentials
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Ping the database to ensure a successful connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	// Create the users table if it doesn't exist
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(50) NOT NULL UNIQUE,
        password TEXT NOT NULL,
        full_name VARCHAR(100) NOT NULL,
        bio TEXT DEFAULT '',
        role ENUM('Writer', 'Admin') DEFAULT 'Writer'
    );`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	log.Println("Connected to the MySQL database and ensured users table exists")

	// Inject the DB connection into the user package
	user.SetDB(db)

	// Routes
	http.HandleFunc("/register", user.RegisterUser)
	http.HandleFunc("/login", user.LoginUser)

	// Handle GET and PUT requests for the /profile route
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			user.ProtectedRoute(user.GetProfile)(w, r)
		case http.MethodPut:
			user.ProtectedRoute(user.UpdateProfile)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/admin", user.ProtectedRoute(AdminOnly)) // Protected admin route
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User Management Service is running"))
	})
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Starting user management service on port 8000...")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// AdminOnly handler to demonstrate role-based access control
// @Summary Admin only access
// @Description Demonstrates role-based access control for admin users.
// @Tags Admin
// @Success 200 {string} string "Welcome, Admin!"
// @Failure 403 {string} string "Access denied: Admins only"
// @Router /admin [get]
func AdminOnly(w http.ResponseWriter, r *http.Request) {
	claims, err := user.ClaimsFromContext(r.Context())
	if err != nil || claims.Role != "Admin" {
		http.Error(w, "Access denied: Admins only", http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome, Admin!"))
}
