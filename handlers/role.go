package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Role struct {
	ID     uint   `json:"id,omitempty" db:"id"`
	UserID uint   `json:"user_id,omitempty" db:"user_id"`
	Type   string `json:"type,omitempty" db:"type"`
}

// CreateRole creates a new role
func CreateRole(c *gin.Context, db *gorm.DB) {
	var role Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user ID from the context (assuming user ID is available)
	userID, _ := c.Get("user_id")
	role.UserID = userID.(uint)

	err := db.Create(&role).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		log.Println("Error executing database query:", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Role created successfully"})
}

// EditRole edits an existing role
func EditRole(c *gin.Context, db *gorm.DB) {
	// Extract roleID from the URL parameters
	roleID, err := strconv.Atoi(c.Param("roleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	// Check if db is nil
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	// Check if the role exists
	var existingRole Role
	err = db.First(&existingRole, roleID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// Extract the new data from the request body
	var request struct {
		UserID uint   `json:"userId"`
		Type   string `json:"type"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the role data
	if request.UserID != 0 {
		existingRole.UserID = request.UserID
	}
	if request.Type != "" {
		existingRole.Type = request.Type
	}

	// Save the updated role
	err = db.Save(&existingRole).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit role"})
		log.Println("Error executing database query:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role edited successfully"})
}

// DeleteRole deletes a role by ID
func DeleteRole(c *gin.Context, db *gorm.DB) {
	// Extract roleID from the URL parameters
	roleID, err := strconv.Atoi(c.Param("roleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	// Check if db is nil
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	// Check if the role exists
	var existingRole Role
	err = db.First(&existingRole, roleID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// Delete the role
	err = db.Delete(&existingRole).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		log.Println("Error executing database query:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// GetAllRoles fetches all roles
func GetAllRoles(c *gin.Context, db *gorm.DB) {
	var roles []Role
	err := db.Find(&roles).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all roles"})
		log.Println("Error executing database query:", err)
		return
	}

	c.JSON(http.StatusOK, roles)
}

// GetRoleByID fetches a role by its ID
func GetRoleByID(c *gin.Context, db *gorm.DB) {
	roleID, err := strconv.Atoi(c.Param("roleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var role Role
	err = db.First(&role, roleID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role"})
		log.Println("Error executing database query:", err)
		return
	}

	c.JSON(http.StatusOK, role)
}
