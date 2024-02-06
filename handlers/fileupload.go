package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const uploadPath = "./uploads/"

// UploadFile handles file uploads
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	// Generate a unique filename
	fileName := filepath.Join(uploadPath, file.Filename)

	// Save the file to the server
	if err := c.SaveUploadedFile(file, fileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

// GetUploadedFiles retrieves a list of uploaded files
func GetUploadedFiles(c *gin.Context) {
	files, err := filepath.Glob(uploadPath + "*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
		return
	}

	c.JSON(http.StatusOK, files)
}
