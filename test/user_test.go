package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/Adnen2/tutorial/firstProject/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	dsn := "host=localhost user=admin password=admin dbname=firstdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	err = db.AutoMigrate(&handlers.User{})
	if err != nil{
		panic("Failed to run migrations: " + err.Error())
	}

	return db
}

func TestRegister(t *testing.T) {
	db := setupTestDB()
	defer db.Migrator().DropTable(&handlers.User{}) // Clean up after the test
	router := gin.New()
	router.POST("/register", func(c *gin.Context) {
		handlers.Register(c, db)
	})

	t.Run("RegisterUser", func(t *testing.T) {
		fmt.Println("Running RegisterUser test")
		requestBody := []byte(`{"username": "testuser", "password": "testpassword"}`)
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var user handlers.User
		err = db.Where("username = ?", "testuser").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, "testuser", user.Username)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("testpassword")))
	})
}

func TestLogin(t *testing.T) {
	db := setupTestDB()
	defer db.Migrator().DropTable(&handlers.User{}) // Clean up after the test
	router := gin.New()
	router.POST("/login", func(c *gin.Context) {
		handlers.Login(c, db)
	})

	t.Run("LoginUser", func(t *testing.T) {
		fmt.Println("Running LoginUser test")
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
		testUser := handlers.User{Username: "testuser", Password: string(hashedPassword)}
		db.Create(&testUser)

		requestBody := []byte(`{"username": "testuser", "password": "testpassword"}`)
		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "access_token")
		assert.Contains(t, w.Body.String(), "refresh_token")
	})
}

// Import necessary packages and modules

func TestProfile(t *testing.T) {
	db := setupTestDB()
	defer db.Migrator().DropTable(&handlers.User{}) // Clean up after the test
	router := gin.Default()
	router.GET("/profile", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.Profile(c, db)
	})

	t.Run("UserProfile", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/profile", nil)
		assert.NoError(t, err)

		// Assume you set the user_id in the context using AuthMiddleware
		ctx := &gin.Context{}
		// You can set "user_id" in the context manually for testing purposes
		ctx.Set("user_id", 1) // Change 123 to a valid user ID

		w := httptest.NewRecorder()

		// Serve HTTP with the authenticated request
		router.ServeHTTP(w, req.WithContext(ctx))

		assert.Equal(t, http.StatusUnauthorized, w.Code) // Update the expected status code

		// Assuming you return an error message in the response JSON
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check if the error message matches the expected response for unauthorized
		assert.Equal(t, "Unauthorized", response["error"])
	})
}
