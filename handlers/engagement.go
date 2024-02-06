package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Engagement struct {
	ID      int    `json:"engagementId,omitempty" db:"id"`
	PostID  int    `json:"postId,omitempty" db:"post_id"`
	UserID  int    `json:"userId,omitempty" db:"user_id"`
	Like    bool   `json:"like,omitempty" db:"like"`
	Comment string `json:"comment,omitempty" db:"comment"`
}

// CreateEngagement creates an engagement for a post
func CreateEngagement(c *gin.Context, db *gorm.DB) {
	var engagement Engagement
	if err := c.ShouldBindJSON(&engagement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user ID from the context (assuming user ID is available in the context)
	userID, _ := c.Get("user_id")
	engagement.UserID = userID.(int)

	err := db.Create(&engagement).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create engagement"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Engagement created successfully"})
}

// UpdateEngagement updates an existing engagement
func UpdateEngagement(c *gin.Context, db *gorm.DB) {
	// Extract engagementId from the URL parameters
	engagementID, err := strconv.Atoi(c.Param("engagementId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid engagement ID"})
		return
	}

	// Check if db is nil
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	// Check if the engagement exists
	var existingEngagement Engagement
	err = db.First(&existingEngagement, engagementID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Engagement not found"})
		return
	}

	// Extract the updated data from the request body
	var request struct {
		Like    bool   `json:"like"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the engagement
	existingEngagement.Like = request.Like
	existingEngagement.Comment = request.Comment

	// Save the updated engagement
	err = db.Save(&existingEngagement).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update engagement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Engagement updated successfully"})
}

// DeleteEngagement deletes an engagement by ID
func DeleteEngagement(c *gin.Context, db *gorm.DB) {
	// Extract engagementId from the URL parameters
	engagementID, err := strconv.Atoi(c.Param("engagementId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid engagement ID"})
		return
	}

	// Check if db is nil
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	// Check if the engagement exists
	var existingEngagement Engagement
	err = db.First(&existingEngagement, engagementID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Engagement not found"})
		return
	}

	// Delete the engagement
	err = db.Delete(&existingEngagement).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete engagement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Engagement deleted successfully"})
}

// GetEngagementsForPost retrieves all engagements for a specific post
func GetEngagementsForPost(c *gin.Context, db *gorm.DB) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var engagements []Engagement
	err = db.Where("post_id = ?", postID).Find(&engagements).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch engagements"})
		return
	}

	c.JSON(http.StatusOK, engagements)
}
