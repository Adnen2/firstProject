package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Company struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Teams       []User `json:"teams" gorm:"foreignKey:CompanyID"`
}

// CreateCompany creates a new company
func CreateCompany(c *gin.Context, db *gorm.DB) {
	var company Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	c.JSON(http.StatusCreated, company)
}

// GetCompanyByID retrieves a company by ID
func GetCompanyByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var company Company
	if err := db.Preload("Teams").First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}
	c.JSON(http.StatusOK, company)
}

// UpdateCompany updates a company
func UpdateCompany(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var company Company
	if err := db.First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	c.JSON(http.StatusOK, company)
}

// DeleteCompany deletes a company
func DeleteCompany(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var company Company
	if err := db.First(&company, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	if err := db.Delete(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}
