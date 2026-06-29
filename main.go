package main

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type File struct {
	ID       string
	Filename string
	URL      string
	Type     string
}

var id string
var out *os.File
var file multipart.File
var header *multipart.FileHeader
var url string

func Service(file multipart.File, header *multipart.FileHeader) (string, error) {
	id = uuid.New().String()
	Filename := id + filepath.Ext(header.Filename)
	path := "upload/" + Filename
	os.Mkdir("upload", 0755)
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}
	url := "http://localhost:8080/upload"
	return url, nil
}
func Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()
	url, err := Service(file, header)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"url": url,
	})
}
func main() {
	r := gin.Default()

	r.POST("/upload", Upload, func(c *gin.Context) {
		c.JSON(200, gin.H{"file go to up": true})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
