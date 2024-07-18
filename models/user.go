package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              int64     `json:"id"`
	Username        string    `json:"username"`
	Password        string    `json:"-"`
	Email           string    `json:"email"`
	RegistrationDate time.Time `json:"registration_date"`
}

var dbpool *pgxpool.Pool

func init() {
	var err error
	dbpool, err = pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to allbase: %v\n", err)
		os.Exit(1)
	}
}

func CreateUser(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	sqlStatement := `INSERT INTO users (username, password, email, registration_date) VALUES ($1, $2, $3, $4) RETURNING id`
	err = dbpool.QueryRow(context.Background(), sqlStatement, user.Username, user.Password, user.Email, user.Registrationdate).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByUsername(username string) (*User, error) {
	user := &User{}

	err := dbpool.QueryRow(context.Background(), `SELECT id, username, password, email, registration_date FROM users WHERE username = $1`, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.RegistrationDate)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func main() {
	newUser := &User{
		Username:        "newUser",
		Password:        "password123",
		Email:           "newuser@example.com",
		RegistrationDate: time.Now(),
	}

	err := CreateUser(newUser)
	if err != nil {
		fmt.Printf("Error creating user: %s\n", err)
		return
	}
	fmt.Printf("User created successfully: %v\n", newUser)

	user, err := GetUserByUsername("newUser")
	if err != nil {
		fmt.Printf("Error fetching user: %s\n", err)
		return
	}
	fmt.Printf("User fetched successfully: %v\n", user)
}