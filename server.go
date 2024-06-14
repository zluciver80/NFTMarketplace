package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

var (
	port  = os.Getenv("PORT")
	users []map[string]string
	nfts  []map[string]string
)

func init() {
	users = append(users, map[string]string{"id": "1", "name": "John Doe"})
	nfts = append(nfts, map[string]string{"id": "1", "name": "CryptoKitty #001", "owner": "John Doe"})
	log.Println("Initial users and NFTs loaded.")
}

func main() {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	authMiddleware := gin.BasicAuth(gin.Accounts{
		"admin": "password",
	})

	nftRoutes := router.Group("/nfts", authMiddleware)
	{
		nftRoutes.GET("/", GetNFTs)
		nftRoutes.POST("/", CreateNFT)
		nftRoutes.PUT("/:id", UpdateNFT)
		nftRoutes.DELETE("/:id", DeleteNFT)
	}

	userRoutes := router.Group("/users", authMiddleware)
	{
		userRoutes.GET("/", GetUsers)
		userRoutes.POST("/", CreateUser)
		userRoutes.PUT("/:id", UpdateUser)
		userRoutes.DELETE("/:id", DeleteUser)
	}

	router.POST("/requests", HandleQueue)

	router.Run(":" + port)
}

func GetNFTs(c *gin.Context) {
	c.JSON(http.StatusOK, nfts)
}

func CreateNFT(c *gin.Context) {
	var nft map[string]string
	if err := c.ShouldBindJSON(&nft); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	nfts = append(nfts, nft)

	log.Println("NFT created:", nft)
	c.JSON(http.StatusCreated, nft)
}

func UpdateNFT(c *gin.Context) {
	id := c.Param("id")
	var nft map[string]string
	if err := c.ShouldBindJSON(&nft); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, item := range nfts {
		if item["id"] == id {
			nfts[i] = nft
			log.Println("NFT updated:", nft)
			c.JSON(http.StatusOK, nft)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "NFT not found"})
}

func DeleteNFT(c *gin.Context) {
	id := c.Param("id")

	for i, item := range nfts {
		if item["id"] == id {
			nfts = append(nfts[:i], nfts[i+1:]...)
			log.Println("NFT deleted:", item)
			c.JSON(http.StatusOK, gin.H{"message": "NFT deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "NFT not found"})
}

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var user map[string]string
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users = append(users, user)

	log.Println("User created:", user)
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user map[string]string
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, item := range users {
		if item["id"] == id {
			users[i] = user
			log.Println("User updated:", user)
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	for i, item := range users {
		if item["id"] == id {
			users = append(users[:i], users[i+1:]...)
			log.Println("User deleted:", item)
			c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not Avaialble"})
}

func HandleQueue(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Handler not implemented yet"})
}