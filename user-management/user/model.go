package user

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// User struct represents a user in the system.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
	Role     string `json:"role"` // Role can be 'Writer' or 'Admin'
}

// JWTClaims defines the claims for the JWT token.
type JWTClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"` // Include user role in the token
	jwt.RegisteredClaims
}

// Secret key used to sign tokens
var jwtKey = []byte("my_secret_key")

// HashPassword generates a bcrypt hash of the password.
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password is correct.
func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}

// GenerateJWT generates a JWT token for the user.
func (u *User) GenerateJWT() (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &JWTClaims{
		Username: u.Username,
		Role:     u.Role, // Add the role to the claims
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
