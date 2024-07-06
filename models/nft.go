package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "os"
    _ "github.com/lib/pq"
)

type NFT struct {
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Creator     string  `json:"creator"`
    Price       float64 `json:"price"`
    UserID      uint    `json:"userID"`
}

var db *sql.DB // Use lowercase for package-private variables

func init() {
    var err error
    connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
    
    // Initialize db connection
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }

    // Configure the database connection pool here if needed
    db.SetMaxOpenConns(25) // Example setting: Adjust based on your usage
    db.SetMaxIdleConns(5)  // Example setting: Keep a few connections idle, ready for use
    db.SetConnMaxLifetime(5 * 60) // Example setting: 5 minutes lifetime for a connection

    // Testing database connectivity
    err = db.Ping()
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
}

func CreateNFT(nft NFT) error {
    // Using context for better control over requests
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    query := `INSERT INTO nfts (title, description, creator, price, user_id) VALUES ($1, $2, $3, $4, $5)`
    _, err := db.ExecContext(ctx, query, nft.Title, nft.Description, nft.Creator, nft.Price, nft.UserID)
    if err != nil {
        log.Printf("Error creating NFT: %v", err)
        return err
    }
    return nil
}

func GetAllNFTs() ([]NFT, error) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    rows, err := db.QueryContext(ctx, "SELECT title, description, creator, price, user_id FROM nfts")
    if err != nil {
        log.Printf("Error retrieving NFTs: %v", err)
        return nil, err
    }
    defer rows.Close() // Ensure rows are closed after function return

    var nfts []NFT
    for rows.Next() {
        var nft NFT
        if err := rows.Scan(&nft.Title, &nft.Description, &nft.Creator, &nft.Price, &nft.UserID); err != nil {
            log.Printf("Error scanning NFT: %v", err)
            return nil, err
        }
        nfts = append(nfts, nft)
    }
    if err = rows.Err(); err != nil {
        log.Printf("Error during rows iteration: %v", err)
        return nil, err
    }
    return nfts, nil
}

func UpdateNFT(nft NFT, id int) error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    query := `UPDATE nfts SET title=$1, description=$2, creator=$3, price=$4 WHERE id=$5`
    _, err := db.ExecContext(ctx, query, nft.Title, nft.Description, nft.Creator, nft.Price, id)
    if err != nil {
        log.Printf("Error updating NFT: %v", err)
        return err
    }
    return nil
}

func DeleteNFT(id int) error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    query := `DELETE FROM nfts WHERE id=$1`
    _, err := db.Execizable"context.WithCancel(context.Background())"(ctx, query, id)
    if err != nil {
        log.Printf("Error deleting NFT: %v", err)
        return err
    }
    return nil
}

func main() {
    // Your main function code
}