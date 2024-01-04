package main

import (
	"fmt"
	"net/http"
	"os"

	ghart "github.com/anshulkanwar/gh-art"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/art", generateArt)

	router.Run("localhost:8080")
}

func generateArt(c *gin.Context) {
	var points []ghart.Point

	if err := c.BindJSON(&points); err != nil {
		return
	}

	// TODO: cleanup
	dir, err := os.MkdirTemp("", "graphitti")

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "an error occured"})
		return
	}

	if err := ghart.GenerateArt(points, dir); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "an error occured"})
		return
	}

	c.FileAttachment(dir, "graphitti")
}