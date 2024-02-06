package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Adnen2/tutorial/firstProject/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB1() *gorm.DB {
	dsn := "host=localhost user=admin password=admin dbname=firstdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	err = db.AutoMigrate(&handlers.Post{})
	if err != nil {
		panic("Failed to run migrations: " + err.Error())
	}

	return db
}

func TestCreatePost(t *testing.T) {
	db := setupTestDB1()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Set("db", db)
		c.Next()
	})
	router.POST("/create-post", func(c *gin.Context) {
		handlers.CreatePost(c, db)
	})

	t.Run("Create a post", func(t *testing.T) {
		payload := []byte(`{"content": "Test post content"}`)
		req, _ := http.NewRequest("POST", "/create-post", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Post created successfully", response["message"])
	})
}

func TestEditPost(t *testing.T) {
	db := setupTestDB1()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Set("db", db)
		c.Next()
	})
	router.PUT("/edit-post/:postId", func(c *gin.Context) {
		handlers.EditPost(c, db)
	})

	t.Run("Edit a post", func(t *testing.T) {
		// Create a test post in the database
		post := handlers.Post{Content: "Test post content"}
		db.Create(&post)

		// Create a request body
		requestBody := []byte(`{"content": "Updated post content"}`)
		req, _ := http.NewRequest("PUT", "/edit-post/"+strconv.Itoa(post.ID), bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder to capture the response
		w := httptest.NewRecorder()

		// Serve the HTTP request
		router.ServeHTTP(w, req)

		// Assert the response status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Assert the response body
		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Post content edited successfully", response["message"])

		// Assert the post is updated in the database
		var updatedPost handlers.Post
		db.First(&updatedPost, post.ID)
		assert.Equal(t, "Updated post content", updatedPost.Content)
	})
}

func TestDeletePost(t *testing.T) {
	db := setupTestDB1()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Set("db", db)
		c.Next()
	})
	router.DELETE("/delete-post/:postId", func(c *gin.Context) {
		handlers.DeletePost(c, db)
	})

	t.Run("Delete a post", func(t *testing.T) {
		// Create a test post in the database
		post := handlers.Post{Content: "Test post content"}
		db.Create(&post)

		// Create a request
		req, _ := http.NewRequest("DELETE", "/delete-post/"+strconv.Itoa(post.ID), nil)

		// Create a response recorder to capture the response
		w := httptest.NewRecorder()

		// Serve the HTTP request
		router.ServeHTTP(w, req)

		// Assert the response status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Assert the response body
		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Post deleted successfully", response["message"])

		// Assert the post is deleted from the database
		var deletedPost handlers.Post
		db.First(&deletedPost, post.ID)
		assert.Equal(t, 0, deletedPost.ID)
	})
}

func TestGetPostByID(t *testing.T) {
	db := setupTestDB1()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Set("db", db)
		c.Next()
	})
	router.GET("/post/:postId", func(c *gin.Context) {
		handlers.GetPostByID(c, db)
	})

	t.Run("Get a post by ID", func(t *testing.T) {
		// Create a test post in the database
		post := handlers.Post{Content: "Test post content"}
		db.Create(&post)

		// Create a request
		req, _ := http.NewRequest("GET", "/post/"+strconv.Itoa(post.ID), nil)

		// Create a response recorder to capture the response
		w := httptest.NewRecorder()

		// Serve the HTTP request
		router.ServeHTTP(w, req)

		// Assert the response status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Assert the response body
		var response map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, post.Content, response["content"])
	})
}

func TestGetAllPosts(t *testing.T) {
	db := setupTestDB1()
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Set("db", db)
		c.Next()
	})
	router.GET("/posts", func(c *gin.Context) {
		handlers.GetAllPosts(c, db)
	})

	t.Run("Get all posts", func(t *testing.T) {
		getReq, _ := http.NewRequest("GET", "/posts", nil)

		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)

		assert.Equal(t, http.StatusOK, getW.Code)
		var getAllResponse []handlers.Post
		_ = json.Unmarshal(getW.Body.Bytes(), &getAllResponse)
		assert.GreaterOrEqual(t, len(getAllResponse), 1)
	})
}
