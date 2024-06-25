package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

var db *sql.DB

func init() {
	var err error
	// Connecting to the database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks the password against the hashed version
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// RegisterUser registers a new user
func RegisterUser(username, email, password string) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	var user User
	err = db.QueryRow("INSERT INTO users(username, email, password) VALUES($1, $2, $3) RETURNING id, username, email",
		username, email, hashedPassword).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(id int) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserDetails updates a user's details
func UpdateUserDetails(id int, username, email string) (*User, error) {
	_, err := db.Exec("UPDATE users SET username = $1, email = $2 WHERE id = $3",
		username, email, id)
	if err != nil {
		return nil, err
	}

	return GetUserByID(id)
}

// AuthenticateUser checks if the user credentials are correct
func AuthenticateUser(email, password string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	if !CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	user.Password = "" // never return the password
	return &user, nil
}

func main() {
	// Example of how to call these functions here. Remember to handle errors in real code.
	fmt.Println("Solana NFT Marketplace Users Management")
}