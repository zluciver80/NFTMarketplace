package main

import (
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

var DB *sql.DB

func init() {
    var err error
    connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
}

func CreateNFT(nft NFT) error {
    query := `INSERT INTO nfts (title, description, creator, price, user_id) VALUES ($1, $2, $3, $4, $5)`
    _, err := DB.Exec(query, nft.Title, nft.Description, nft.Creator, nft.Price, nft.UserID)
    if err != nil {
        log.Printf("Error creating NFT: %v", err)
        return err
    }
    return nil
}

func GetAllNFTs() ([]NFT, error) {
    rows, err := DB.Query("SELECT title, description, creator, price, user_id FROM nfts")
    if err != nil {
        log.Printf("Error retrieving NFTs: %v", err)
        return nil, err
    }
    defer rows.Close()

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
    query := `UPDATE nfts SET title=$1, description=$2, creator=$3, price=$4 WHERE id=$5`
    _, err := DB.Exec(query, nft.Title, nft.Description, nft.Creator, nft.Price, id)
    if err != nil {
        log.Printf("Error updating NFT: %v", err)
        return err
    }
    return nil
}

func DeleteNFT(id int) error {
    query := `DELETE FROM nfts WHERE id=$1`
    _, err := DB.Exec(query, id)
    if err != nil {
        log.Printf("Error deleting NFT: %v", err)
        return err
    }
    return nil
}

func main() {
    // Your main function code
}