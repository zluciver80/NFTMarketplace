package main

import (
    "context"
    "fmt"
    "os"
    "sync"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type NFT struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Name        string             `bson:"name"`
    Description string             `bson:"description"`
    ImageURL    string             `bson:"imageURL"`
    Owner       string             `bson:"owner"`
    CreatedAt   time.Time          `bson:"createdAt"`
}

var (
    collection *mongo.Collection
    ctx        context.Context
    cache      []*NFT
    mu         sync.RWMutex
)

func initDB() {
    mongoURI := os.Getenv("MONGO_URI")
    clientOptions := options.Client().ApplyURI(mongoURI).SetMaxPoolSize(50)

    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        panic(err)
    }

    ctx = context.Background()
    collection = client.Database("SolanaNFTMarketplace").Collection("nfts")
}

func CreateNFTs(nfts []NFT) (*mongo.InsertManyResult, error) {
    docs := make([]interface{}, len(nfts))
    for i, nft := range nfts {
        nft.ID = primitive.NewObjectID()
        nft.CreatedAt = time.Now()
        docs[i] = nft
    }

    result, err := collection.InsertMany(ctx, docs)
    if err != nil {
        return nil, err
    }

    invalidateCache()
    return result, nil
}

func GetNFTs() ([]*NFT, error) {
    mu.RLock()
    if cache != nil {
        mu.RUnlock()
        return cache, nil
    }
    mu.RUnlock()

    var nfts []*NFT

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var nft NFT
        err := cursor.Decode(&nft)
        if err != nil {
            return nil, err
        }

        nfts = append(nfts, &nft)
    }

    mu.Lock()
    cache = nfts
    mu.Unlock()

    return nfts, nil
}

func UpdateNFT(id string, update bson.M) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    filter := bson.M{"_id": objID}
    _, err = collection.UpdateOne(ctx, filter, bson.M{
        "$set": update,
    })

    invalidateCache()
    return err
}

func DeleteNFT(id string) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    _, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil {
        return err
    }

    invalidateCache()
    return nil
}

func invalidateCache() {
    mu.Lock()
    cache = nil
    mu.Unlock()
}

func main() {
    initDB()

    nfts := []NFT{
        {
            Name:        "Example NFT 1",
            Description: "This is an example NFT 1",
            ImageURL:    "https://example.com/nft1.jpg",
            Owner:       "John Doe",
        },
        {
            Name:        "Example NFT 2",
            Description: "This is an example NFT 2",
            ImageURL:    "https://example.com/nft2.jpg",
            Owner:       "Jane Doe",
        },
    }

    result, err := CreateNFTs(nfts)
    if err != nil {
        fmt.Println("Error creating NFTs:", err)
        return
    }
    fmt.Println("NFTs created:", result.InsertedIDs)

    existingNfts, err := GetNFTs()
    if err != nil {
        fmt.Println("Error getting NFTs:", err)
        return
    }
    for _, n := range existingNfts {
        fmt.Printf("NFT: %#v\n", n)
    }

    if len(existingNfts) > 0 {
        err = UpdateNFT(existingNfts[0].ID.Hex(), bson.M{"name": "Updated NFT Name"})
        if err != nil {
            fmt.Println("Error updating NFT:", err)
            return
        }
        fmt.Println("NFT updated")

        err = DeleteNFT(existingNfts[0].ID.Hex())
        if err != nil {
            fmt.Println("Error deleting NFT:", err)
            return
        }
        fmt.Println("NFT deleted")
    }
}