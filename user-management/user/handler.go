package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Define a custom key type to prevent collisions in context keys
type key int

const (
	claimsKey key = iota
)

var db *sql.DB

// SetDB sets the database connection for the user package
func SetDB(database *sql.DB) {
	db = database
}

// ContextWithClaims adds JWT claims to the request context
func ContextWithClaims(ctx context.Context, claims *JWTClaims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// ClaimsFromContext retrieves the JWT claims from the request context
func ClaimsFromContext(ctx context.Context) (*JWTClaims, error) {
	claims, ok := ctx.Value(claimsKey).(*JWTClaims)
	if !ok {
		return nil, errors.New("no claims in context")
	}
	return claims, nil
}

// RegistrationRequest represents the structure for user registration input
type RegistrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// LoginRequest represents the structure for login input
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ProfileRequest represents the structure for profile update requests
type ProfileRequest struct {
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
}

// RegisterUser handles the registration of a new user.
// @Summary Register a new user
// @Description Register a new user by providing username, password, and full name
// @Tags User
// @Accept  json
// @Produce  json
// @Param   user  body  RegistrationRequest  true  "User Registration"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
// RegisterUser handles the registration of a new user.
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hash the password
	user := User{
		Username: req.Username,
		FullName: req.FullName,
		Role:     "Writer", // Default role
	}
	if err := user.HashPassword(req.Password); err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert the new user into the database
	insertQuery := `INSERT INTO users (username, password, full_name, bio, role) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(insertQuery, req.Username, user.Password, user.FullName, user.Bio, user.Role)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	log.Printf("User registered: %s", req.Username)

	// Return success message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

// LoginUser handles user login and returns a JWT token.
// @Summary Login a user
// @Description Logs in a user and returns a JWT token
// @Tags User
// @Accept  json
// @Produce  json
// @Param   login  body  LoginRequest  true  "User Login"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
// LoginUser handles user login and returns a JWT token.
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if the user exists in the database
	var user User
	query := `SELECT id, username, password, full_name, bio, role FROM users WHERE username = ?`
	err = db.QueryRow(query, req.Username).Scan(&user.ID, &user.Username, &user.Password, &user.FullName, &user.Bio, &user.Role)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Check if the password matches
	if err := user.CheckPassword(req.Password); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token with role
	token, err := user.GenerateJWT()
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return the token as a response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// GetProfile retrieves the profile of the currently logged-in user.
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags Profile
// @Produce  json
// @Success 200 {object} User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile [get]
func GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := ClaimsFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user User
	query := `SELECT id, username, full_name, bio, role FROM users WHERE username = ?`
	err = db.QueryRow(query, claims.Username).Scan(&user.ID, &user.Username, &user.FullName, &user.Bio, &user.Role)
	if err == sql.ErrNoRows {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Return the user profile as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateProfile allows users to update their profile.
// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param   profile  body  ProfileRequest  true  "User Profile"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile [put]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := ClaimsFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ProfileRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Update the user's profile in the database
	updateQuery := `UPDATE users SET full_name = ?, bio = ? WHERE username = ?`
	_, err = db.Exec(updateQuery, req.FullName, req.Bio, claims.Username)
	if err != nil {
		log.Printf("Failed to update profile: %v", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	log.Printf("Profile updated for user: %s", claims.Username)

	// Return success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile updated successfully",
	})
}

// Middleware for protected routes
func ProtectedRoute(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		r = r.WithContext(ContextWithClaims(r.Context(), claims))
		next(w, r)
	}
}
