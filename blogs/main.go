package main

import (
	"blogs/blog"
	_ "blogs/docs" // For Swagger documentation
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	httpSwagger "github.com/swaggo/http-swagger"
)

var db *sql.DB

// ProtectedRoute is a simple middleware that ensures authentication.
func ProtectedRoute(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// In a real scenario, you would check for JWT or session authentication here
		// If not authenticated, return Unauthorized
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Pass to the next handler if authenticated
		next(w, r)
	}
}

// @title Blog Management API
// @version 1.0
// @description API for handling blog operations (CRUD) with role-based access control.
// @host localhost:8001
// @BasePath /

func main() {
	// MySQL connection string (ensure you use your credentials)
	dsn := "root:@tcp(127.0.0.1:3306)/blog_management"
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

	// Create the blogs table if it doesn't exist
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS blogs (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        author VARCHAR(100) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create blogs table: %v", err)
	}

	log.Println("Connected to the MySQL database and ensured blogs table exists")

	// Inject the DB connection into the blog package
	blog.SetDB(db)

	// Routes
	http.HandleFunc("/blogs", ProtectedRoute(blog.GetBlogs))          // GET all blogs
	http.HandleFunc("/blogs/create", ProtectedRoute(blog.CreateBlog)) // POST a new blog
	http.HandleFunc("/blogs/update", ProtectedRoute(blog.UpdateBlog)) // PUT update a blog
	http.HandleFunc("/blogs/delete", ProtectedRoute(blog.DeleteBlog)) // DELETE a blog
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Starting blog management service on port 8001...")
	err = http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
