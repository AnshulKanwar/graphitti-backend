package main

import (
	"fmt"
	"net/http"
	"os"
	"path"

	ghart "github.com/anshulkanwar/gh-art"
	"github.com/anshulkanwar/graphitti-backend/internal"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.POST("/generate", generate)
	router.GET("/download/:filename", downloadZip)

	router.Run()
}

func generate(c *gin.Context) {
	var data struct {
		Config ghart.Config
		Art []ghart.Point
	}

	if err := c.BindJSON(&data); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "an error occured"})
		return
	}

	// TODO: cleanup
	// TODO: better error handling
	dir, err := os.MkdirTemp("", "graphitti")

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "an error occured"})
		return
	}

	if err := ghart.GenerateArt(data.Art, data.Config, dir); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "an error occured"})
		return
	}

	zip, err := os.CreateTemp("", "graphitti-*.zip")

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "an error occured"})
		return
	}

	defer zip.Close()

	if err := internal.Zip(zip, dir); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "an error occured"})
		return
	}

	_, filename := path.Split(zip.Name())

	c.IndentedJSON(http.StatusOK, gin.H{"filename": filename})
}

func downloadZip(c *gin.Context) {
	filename := c.Param("filename")

	// TODO: allow only graphitti zips file to be downloaded
	path := path.Join(os.TempDir(), filename)

	c.FileAttachment(path, "graphitti.zip")
}
