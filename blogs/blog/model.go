package blog

import (
	"database/sql"
	"time"
)

// Blog represents a blog post.
type Blog struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

// SetDB sets the database connection for the blog package.
func SetDB(database *sql.DB) {
	db = database
}

// CreateBlog inserts a new blog post into the database.
func (b *Blog) CreateBlog() error {
	query := `INSERT INTO blogs (title, content, author) VALUES (?, ?, ?)`
	_, err := db.Exec(query, b.Title, b.Content, b.Author)
	return err
}

// GetAllBlogs retrieves all blog posts from the database.
func GetAllBlogs() ([]Blog, error) {
	query := `SELECT id, title, content, author, created_at FROM blogs`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []Blog
	for rows.Next() {
		var blog Blog
		if err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Author, &blog.CreatedAt); err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

// GetBlogByID retrieves a single blog post by its ID.
func GetBlogByID(id int) (*Blog, error) {
	query := `SELECT id, title, content, author, created_at FROM blogs WHERE id = ?`
	row := db.QueryRow(query, id)

	var blog Blog
	err := row.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Author, &blog.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &blog, nil
}

// UpdateBlog updates an existing blog post in the database.
func (b *Blog) UpdateBlog() error {
	query := `UPDATE blogs SET title = ?, content = ? WHERE id = ?`
	_, err := db.Exec(query, b.Title, b.Content, b.ID)
	return err
}

// DeleteBlogFromModel deletes a blog post from the database.
func DeleteBlogFromModel(id int) error {
	query := `DELETE FROM blogs WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}
