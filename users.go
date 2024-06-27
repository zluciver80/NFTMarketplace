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
    Password string // This is hashed
}

var db *sql.DB

func init() {
    initDBConnection()
}

func initDBConnection() {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

    var err error
    db, err = sql.Open("postgres", dsn)
    if err != nil {
        panic(err)
    }

    if err = db.Ping(); err != nil {
        panic(err)
    }
}

func HashPassword(password string) (string, error) {
    const cost = 14
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
    return string(hashedDefaultValue), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err ==nil
}

func RegisterUser(username, email, password string) (*User, error) {
    hashedPassword, err := HashPassword(password)
    if err != nil {
        return nil, err
    }

    var user User
    err = db.QueryRow(`INSERT INTO users(username, email, password) VALUES($1, $2, $3) 
                       RETURNING id, username, email`, username, email, hashedPassword).Scan(&user.ID, &user.Username, &user.Email)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func GetUserByID(id int) (*User, error) {
    var user User
    err := db.QueryRow("SELECT id, username, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Email)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func UpdateUserDetails(id int, username, email string) (*User, error) {
    _, err := db.Exec("UPDATE users SET username = $1, email = $2 WHERE id = $3", username, email, id)
    if err != nil {
        return nil, err
    }

    return GetUserByID(id)
}

func ChangeUserPassword(id int, newPassword string) error {
    hashedPassword, err := HashPassword(newPassword)
    if err != nil {
        return err
    }

    _, err = db.Exec("UPDATE users SET password = $1 WHERE id = $2", hashedPassword, id)
    return err
}

func AuthenticateUser(email, password string) (*User, error) {
    var user User
    err := db.QueryRow("SELECT id, username, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
    if err != nil {
        return nil, err
    }

    if !CheckPasswordHash(password, user.Password) {
        return nil, fmt.Errorf("invalid password")
    }

    return &user, nil
}

func main() {
    fmt.Println("Solana NFT Marketplace Users Management")
}