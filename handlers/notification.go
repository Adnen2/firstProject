package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Notification struct {
	ID        int       `json:"notificationId,omitempty" db:"id"`
	UserID    int       `json:"userId,omitempty" db:"user_id"`
	Message   string    `json:"message,omitempty" db:"message"`
	IsRead    bool      `json:"isRead,omitempty" db:"is_read"`
	CreatedAt time.Time `json:"createdAt,omitempty" db:"created_at"`
}

// CreateNotification sends a notification to a user
func CreateNotification(c *gin.Context, db *gorm.DB) {
	var notification Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.Create(&notification).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Notification sent successfully"})
}

// GetNotifications retrieves notifications for a user
func GetNotifications(c *gin.Context, db *gorm.DB) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var notifications []Notification
	err := db.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkNotificationAsRead marks a notification as read
func MarkNotificationAsRead(c *gin.Context, db *gorm.DB) {
	notificationID := c.Param("notificationId")
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var notification Notification
	err := db.First(&notification, notificationID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	if notification.UserID != userID.(int) {
		fmt.Println(userID.(int))
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to mark this notification as read"})
		return
	}

	// Mark the notification as read
	err = db.Model(&notification).Update("is_read", true).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read successfully"})
}
