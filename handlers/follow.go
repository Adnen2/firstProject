package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Follow struct {
	ID          int `json:"followId,omitempty" db:"id"`
	FollowerID  int `json:"followerId,omitempty" db:"follower_id"`
	FollowingID int `json:"followingId,omitempty" db:"following_id"`
}

// FollowUser allows a user to follow another user
func FollowUser(c *gin.Context, db *gorm.DB) {
	var follow Follow
	if err := c.ShouldBindJSON(&follow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set follower ID from the context (assuming user ID is available in the context)
	followerID, _ := c.Get("user_id")
	follow.FollowerID = followerID.(int)

	err := db.Create(&follow).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User followed successfully"})
}

// UnfollowUser allows a user to unfollow another user
func UnfollowUser(c *gin.Context, db *gorm.DB) {
	// Extract following ID from the URL parameters
	followingID, err := strconv.Atoi(c.Param("followingId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid following ID"})
		return
	}

	// Set follower ID from the context (assuming user ID is available in the context)
	followerID, _ := c.Get("user_id")

	// Check if the follow relationship exists
	var existingFollow Follow
	err = db.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existingFollow).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follow relationship not found"})
		return
	}

	// Delete the follow relationship
	err = db.Delete(&existingFollow).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unfollowed successfully"})
}

// GetFollowers retrieves followers for a user
func GetFollowers(c *gin.Context, db *gorm.DB) {
	// Extract user ID from the URL parameters
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var followers []Follow
	err = db.Where("following_id = ?", userID).Find(&followers).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followers"})
		return
	}

	c.JSON(http.StatusOK, followers)
}

// GetFollowings retrieves users that a user is following
func GetFollowings(c *gin.Context, db *gorm.DB) {
	// Extract user ID from the URL parameters
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var followings []Follow
	err = db.Where("follower_id = ?", userID).Find(&followings).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followings"})
		return
	}

	c.JSON(http.StatusOK, followings)
}
