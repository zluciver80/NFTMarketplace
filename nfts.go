package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/momo-driver/bson/primitive"
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

var collection *mongo.Collection

func initDB() {
	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	collection = client.Database("SolanaNFTMarketplace").Collection("nfts")
}

func CreateNFT(nft NFT) (*mongo.InsertOneResult, error) {
	nft.ID = primitive.NewObjectID()
	nft.CreatedAt = time.Now()

	result, err := collection.InsertOne(context.TODO(), nft)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetNFTs() ([]*NFT, error) {
	var nfts []*NFT

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var nft NFT
		err := cursor.Decode(&nft)
		if err != nil {
			return nil, err
		}

		nfts = append(nfts, &nft)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(context.TODO())

	return nfts, nil
}

func UpdateNFT(id string, update bson.M) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{
		"$set": update,
	})

	return err
}

func DeleteNFT(id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != null {
		return err
	}

	return nil
}

func main() {
	initDB()

	nft := NFT{
		Name:        "Example NFT",
		Description: "This is an example NFT",
		ImageURL:    "https://example.com/nft.jpg",
		Owner:       "John Doe",
	}

	result, err := CreateNFT(nft)
	if err != nil {
		fmt.Println("Error creating NFT:", err)
		return
	}
	fmt.Println("NFT created:", result.InsertedID)

	nfts, err := GetNFTs()
	if err != nil {
		fmt.Println("Error getting NFTs:", err)
		return
	}
	for _, n := range nfts {
		fmt.Printf("NFT: %#v\n", n)
	}

	err = UpdateNFT(nfts[0].ID.Hex(), bson.M{"name": "Updated NFT Name"})
	if err != nil {
		fmt.Println("Error updating NFT:", err)
		return
	}
	fmt.Println("NFT updated")

	err = DeleteNFT(nfts[0].ID.Hex())
	if err != nil {
		fmt.Println("Error deleting NFT:", err)
		return
	}
	fmt.Println("NFT deleted")
}