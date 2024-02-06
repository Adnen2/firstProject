package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Post struct {
	ID           int       `json:"postId,omitempty" db:"id"`
	Content      string    `json:"content,omitempty" db:"content"`
	ScheduleTime time.Time `json:"scheduleTime,omitempty" db:"schedule_time"`
	UserID       int       `json:"userId,omitempty" db:"user_id"`
}

func CreatePost(c *gin.Context, db *gorm.DB) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user ID from the context (assuming user ID is available in the context)
	userID, _ := c.Get("user_id")
	post.UserID = userID.(int)

	err := db.Create(&post).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		log.Println("Error executing database query:", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}
func EditPost(c *gin.Context, db *gorm.DB) {
	// Extract postId from the URL parameters
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Check if db is nil
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	// Check if the post exists
	var existingPost Post
	err = db.First(&existingPost, postID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Extract the new content from the request body
	var request struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the post content
	existingPost.Content = request.Content

	// Save the updated post
	err = db.Save(&existingPost).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit post"})
		log.Println("Error executing database query:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post content edited successfully"})
}

// DeletePost deletes a post by ID
func DeletePost(c *gin.Context, db *gorm.DB) {
	// Extract postId from the URL parameters
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Check if db is nil
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	// Check if the post exists
	var existingPost Post
	err = db.First(&existingPost, postID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Delete the post
	err = db.Delete(&existingPost).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		log.Println("Error executing database query:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// GetPostByID gets a post by ID
func GetPostByID(c *gin.Context, db *gorm.DB) {
	// Extract postId from the URL parameters
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Check if db is nil
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	// Check if the post exists
	var post Post
	err = db.First(&post, postID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// GetAllPosts gets all posts
func GetAllPosts(c *gin.Context, db *gorm.DB) {
	// Check if db is nil
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	// Get all posts
	var posts []Post
	err := db.Find(&posts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		log.Println("Error executing database query:", err)
		return
	}

	c.JSON(http.StatusOK, posts)
}
