package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/jackc/pgconn"
    "github.com/jackc/pgx/v4"
    "github.com/jackc/pgx/v4/pgxpool"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID               int64     `json:"id"`
    Username         string    `json:"username"`
    Password         string    `json:"-"`
    Email            string    `json:"email"`
    RegistrationDate time.Time `json:"registration_date"`
}

var dbpool *pgxpool.Pool

func init() {
    var err error
    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        fmt.Fprintln(os.Stderr, "DATABASE_URL is not set")
        os.Exit(1)
    }

    dbpool, err = pgxpool.Connect(context.Background(), databaseURL)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
        os.Exit(1)
    }
}

func CreateUser(user *User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("error hashing password: %w", err)
    }
    user.Password = string(hashedPassword)

    sqlStatement := `INSERT INTO users (username, password, email, registration_date) VALUES ($1, $2, $3, $4) RETURNING id`
    if err := dbpool.QueryRow(context.Background(), sqlStatement, user.Username, user.Password, user.Email, user.RegistrationDate).Scan(&user.ID); err != nil {
        var pgErr *pgconn.PgError
        if ok := pgconn.As(err, &pgErr); ok {
            return fmt.Errorf("database error %d: %s", pgErr.Code, pgErr.Message)
        }
        return fmt.Errorf("error inserting user into database: %w", err)
    }

    return nil
}

func GetUserByUsername(username string) (*User, error) {
    user := &User{}

    err := dbpool.QueryRow(context.Background(), `SELECT id, username, password, email, registrationation_date FROM users WHERE username = $1`, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.RegistrationDate)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("error retrieving user by username: %w", err)
    }

    return user, nil
}

func main() {
    newUser := &User{
        Username:         "newUser",
        Password:         "password123",
        Email:            "newuser@example.com",
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