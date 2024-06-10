package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

var (
	port = os.Getenv("PORT")
)

func main() {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/nfts", GetNFTs)
	router.POST("/nfts", CreateNFT)
	router.PUT("/nfts/:id", UpdateNFT)
	router.DELETE("/nfts/:id", DeleteNFT)

	router.GET("/users", GetUsers)
	router.POST("/users", CreateUser)
	router.PUT("/users/:id", UpdateUser)
	router.DELETE("/users/:id", DeleteUser)

	router.POST("/requests", HandleRequest)

	router.Run(":" + port)
}

func GetNFTs(c *gin.Context) {
}

func CreateNFT(c *gin.Context) {
}

func UpdateNFT(c *gin.Context) {
}

func DeleteNFT(c *gin.Context) {
}

func GetUsers(c *gin.Context) {
}

func CreateUser(c *gin.Context) {
}

func UpdateUser(c *gin.Context) {
}

func DeleteUser(c *gin.Uri) {
}

func HandleQueue(c *gin.Context) {
}