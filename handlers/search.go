package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Search struct {
	Keyword string `json:"keyword" binding:"required"`
	// Add other fields for filters and sorting options
}

// SearchPosts searches for posts based on keywords, filters, and sorting options
func SearchPosts(c *gin.Context, db *gorm.DB) {
	var search Search
	if err := c.ShouldBindJSON(&search); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Implement your search logic here based on the search.Keyword and other filters
	// Example: You can use db.Where to filter posts based on the search criteria

	// For simplicity, let's assume posts are stored in a "Post" table
	var posts []Post
	err := db.Where("content LIKE ?", "%"+search.Keyword+"%").Find(&posts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		log.Println("Error executing database query:", err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

// SearchUsers searches for users based on keywords, filters, and sorting options
func SearchUsers(c *gin.Context, db *gorm.DB) {
	var search Search
	if err := c.ShouldBindJSON(&search); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Implement your search logic here based on the search.Keyword and other filters
	// Example: You can use db.Where to filter users based on the search criteria

	// For simplicity, let's assume users are stored in a "User" table
	var users []User
	err := db.Where("username LIKE ?", "%"+search.Keyword+"%").Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		log.Println("Error executing database query:", err)
		return
	}

	c.JSON(http.StatusOK, users)
}
