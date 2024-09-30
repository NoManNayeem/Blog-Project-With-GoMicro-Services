package blog

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims defines the claims structure for the JWT token.
type JWTClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type key int

const (
	claimsKey key = iota
)

// ClaimsFromContext retrieves the JWT claims from the request context.
func ClaimsFromContext(ctx context.Context) (*JWTClaims, error) {
	claims, ok := ctx.Value(claimsKey).(*JWTClaims)
	if !ok {
		return nil, errors.New("no claims in context")
	}
	return claims, nil
}

// CreateBlog handles the creation of a new blog post.
// @Summary Create a new blog post
// @Description Allows a writer to create a new blog post
// @Tags Blog
// @Accept  json
// @Produce  json
// @Param   blog  body  Blog  true  "Blog Post"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /blogs/create [post]
func CreateBlog(w http.ResponseWriter, r *http.Request) {
	var blog Blog
	err := json.NewDecoder(r.Body).Decode(&blog)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	claims, err := ClaimsFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	blog.Author = claims.Username

	err = blog.CreateBlog()
	if err != nil {
		log.Printf("Failed to create blog: %v", err)
		http.Error(w, "Failed to create blog", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Blog post created successfully",
	})
}

// GetBlogs retrieves all blog posts.
// @Summary Get all blog posts
// @Description Retrieves all blog posts from the database
// @Tags Blog
// @Produce  json
// @Success 200 {array} Blog
// @Failure 500 {object} map[string]string
// @Router /blogs [get]
func GetBlogs(w http.ResponseWriter, r *http.Request) {
	blogs, err := GetAllBlogs()
	if err != nil {
		log.Printf("Failed to retrieve blogs: %v", err)
		http.Error(w, "Failed to retrieve blogs", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(blogs)
}

// UpdateBlog handles the update of an existing blog post.
// @Summary Update a blog post
// @Description Allows a writer to update their blog post
// @Tags Blog
// @Accept  json
// @Produce  json
// @Param   blog  body  Blog  true  "Updated Blog Post"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /blogs/update [put]
func UpdateBlog(w http.ResponseWriter, r *http.Request) {
	var blog Blog
	err := json.NewDecoder(r.Body).Decode(&blog)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	existingBlog, err := GetBlogByID(blog.ID)
	if err != nil || existingBlog == nil {
		http.Error(w, "Blog post not found", http.StatusNotFound)
		return
	}

	claims, err := ClaimsFromContext(r.Context())
	if err != nil || claims.Username != existingBlog.Author {
		http.Error(w, "Forbidden: You can only update your own blog post", http.StatusForbidden)
		return
	}

	err = blog.UpdateBlog()
	if err != nil {
		log.Printf("Failed to update blog: %v", err)
		http.Error(w, "Failed to update blog", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Blog post updated successfully",
	})
}

// DeleteBlog handles the deletion of a blog post.
// @Summary Delete a blog post
// @Description Allows a writer to delete their blog post
// @Tags Blog
// @Produce  json
// @Param   id  query  int  true  "Blog ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /blogs/delete [delete]
func DeleteBlog(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing blog ID", http.StatusBadRequest)
		return
	}
	blogID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	blog, err := GetBlogByID(blogID)
	if err != nil || blog == nil {
		http.Error(w, "Blog post not found", http.StatusNotFound)
		return
	}

	claims, err := ClaimsFromContext(r.Context())
	if err != nil || claims.Username != blog.Author {
		http.Error(w, "Forbidden: You can only delete your own blog post", http.StatusForbidden)
		return
	}

	// Call the model's DeleteBlog function
	err = DeleteBlogFromModel(blogID)
	if err != nil {
		log.Printf("Failed to delete blog: %v", err)
		http.Error(w, "Failed to delete blog", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Blog post deleted successfully",
	})
}
