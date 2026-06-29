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
	Filename string
	URL      string
}

var files = make(map[string]File)

func Service(file multipart.File, header *multipart.FileHeader) (File, error) {
	id := uuid.New().String()
	filename := id + filepath.Ext(header.Filename)
	path := "upload/" + filename
	_ = os.Mkdir("upload", 0755)
	out, err := os.Create(path)
	if err != nil {
		return File{}, err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return File{}, err
	}
	f := File{
		Filename: filename,
		URL:      "/upload/" + filename,
	}
	return f, nil
}
func Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	f, err := Service(file, header)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"file": f,
	})
}

func main() {
	r := gin.Default()

	r.POST("/upload", Upload)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Static("/upload", "./upload")
	r.Run(":" + port)
}
