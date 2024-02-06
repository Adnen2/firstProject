package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PostView represents a post view entity
type PostView struct {
	ID        int       `json:"viewId,omitempty" db:"id"`
	PostID    int       `json:"postId,omitempty" db:"post_id"`
	UserID    int       `json:"userId,omitempty" db:"user_id"`
	Timestamp time.Time `json:"timestamp,omitempty" db:"timestamp"`
}

// EngagementMetrics represents engagement metrics for a post
type EngagementMetrics struct {
	PostID  int   `json:"postId,omitempty" db:"post_id"`
	Like    int64 `json:"like,omitempty" db:"like"`
	Comment int64 `json:"comment,omitempty" db:"comment"`
	View    int64 `json:"view,omitempty" db:"view"`
	UserID  int   `json:"userId,omitempty" db:"user_id"`
}

// TrackPostView tracks a view for a post
func TrackPostView(c *gin.Context, db *gorm.DB) {
	var postView PostView
	if err := c.ShouldBindJSON(&postView); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the post with the specified ID exists
	var existingPost Post
	err := db.First(&existingPost, postView.PostID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Set user ID from the context (assuming user ID is available in the context)
	userID, _ := c.Get("user_id")
	postView.UserID = userID.(int)
	postView.Timestamp = time.Now()

	err = db.Create(&postView).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track post view"})
		log.Println("Error executing database query:", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post view and engagement metrics tracked successfully"})
}

// GetPostAnalytics retrieves analytics data for a post
func GetPostAnalytics(c *gin.Context, db *gorm.DB) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Assuming "likes" and "comments" are boolean fields in the "Engagement" table
	var engagementMetrics EngagementMetrics

	// Get the count of likes and comments for the specified PostID
	err = db.Model(&Engagement{}).
		Where("post_id = ? AND \"like\" = true", postID). // Corrected: use backticks around like
		Count(&engagementMetrics.Like).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve likes count"})
		log.Println("Error executing database query:", err)
		return
	}

	err = db.Model(&Engagement{}).
		Where("post_id = ? AND comment IS NOT NULL", postID).
		Count(&engagementMetrics.Comment).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments count"})
		log.Println("Error executing database query:", err)
		return
	}
	err = db.Model(&PostView{}).
		Where("post_id = ?", postID).
		Count(&engagementMetrics.View).
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments View"})
		log.Println("Error executing database query:", err)
		return
	}

	// Set other fields in the EngagementMetrics structure
	engagementMetrics.PostID = int(postID)

	// Set user ID from the context (assuming user ID is available in the context)
	userID, _ := c.Get("user_id")
	engagementMetrics.UserID = userID.(int)
	err = db.Create(&engagementMetrics).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save analyse"})
		log.Println("Error executing database query:", err)
		return
	}
	c.JSON(http.StatusOK, engagementMetrics)
}
