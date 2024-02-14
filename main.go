package main

import (
	"fmt"
	"log"
	"net/http"

	handlers "github.com/Adnen2/tutorial/firstProject/handlers"
	"github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbName = "firstdb1"
)

func initDB() {
	var err error

	// Update these values with your actual database credentials and settings
	username := "admin"
	password := "admin"

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable", username, password, dbName)
	// Assign the connection to the global db variable, not creating a new local one
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Auto-migrate the User model
	err = db.AutoMigrate(&handlers.User{})
	// Auto-migrate the Post model
	err1 := db.AutoMigrate(&handlers.Post{})
	// Auto-migrate the Engagement model
	err2 := db.AutoMigrate(&handlers.Engagement{})
	// Auto-migrate the Notification model
	err3 := db.AutoMigrate(&handlers.Notification{})
	// Auto-migrate the Follow model
	err4 := db.AutoMigrate(&handlers.Follow{})
	// Auto-migrate the Search model
	err5 := db.AutoMigrate(&handlers.Search{})
	// Auto-migrate the PostView model
	err6 := db.AutoMigrate(&handlers.PostView{})
	// Auto-migrate the EngagementMetrics model
	err7 := db.AutoMigrate(&handlers.EngagementMetrics{})
	// Auto-migrate the Company model
	err8 := db.AutoMigrate(&handlers.Company{})
	// Auto-migrate the Role model
	err9 := db.AutoMigrate(&handlers.Role{})
	if err != nil && err1 != nil && err2 != nil && err3 != nil && err4 != nil && err5 != nil && err6 != nil && err7 != nil && err8 != nil && err9 != nil {
		log.Fatal("Error auto-migrating database:", err)
	}

	// This line is optional but closes the database connection when the main function exits
	//  sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error getting underlying *sql.DB:", err)
	}
	// defer sqlDB.Close()

	fmt.Println("Connected to PostgreSQL and auto-migrated tables!")
}

func main() {
	initDB()

	router := gin.Default()
	// user routes
	router.POST("/register", func(c *gin.Context) {
		handlers.Register(c, db)
	})
	router.POST("/login", func(c *gin.Context) {
		handlers.Login(c, db)
	})
	router.GET("/profile", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.Profile(c, db)
	})

	router.PUT("/update-profile", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.UpdateProfile(c, db)
	})
	//Router of post
	router.POST("/create-post", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.CreatePost(c, db)
	})

	router.PUT("/edit-post/:postId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.EditPost(c, db)
	})
	router.DELETE("/posts/:postId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.DeletePost(c, db)
	})
	router.GET("/posts/:postId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.GetPostByID(c, db)
	})
	router.GET("/posts", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.GetAllPosts(c, db)
	})
	//Engagement router
	router.POST("/engagements", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.CreateEngagement(c, db)
	})
	router.PUT("/engagements/:engagementId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.UpdateEngagement(c, db)
	})
	router.DELETE("/engagements/:engagementId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.DeleteEngagement(c, db)
	})
	router.GET("/engagements/:postId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.GetEngagementsForPost(c, db)
	})
	// Notification routes
	router.POST("/notifications", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.CreateNotification(c, db)
	})
	router.GET("/notifications", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.GetNotifications(c, db)
	})
	router.PATCH("/notifications/:notificationId/read", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.MarkNotificationAsRead(c, db)
	})
	// Follow/Unfollow routes
	router.POST("/follow", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.FollowUser(c, db)
	})

	router.DELETE("/unfollow/:followingId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.UnfollowUser(c, db)
	})

	router.GET("/followers/:userId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.GetFollowers(c, db)
	})

	router.GET("/followings/:userId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.GetFollowings(c, db)
	})
	// Add these search routes
	router.POST("/search/posts", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.SearchPosts(c, db)
	})

	router.POST("/search/users", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.SearchUsers(c, db)
	})
	// New routes for Analytics
	router.POST("/track-post-view", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.TrackPostView(c, db)
	})
	router.GET("/post-analytics/:postId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.GetPostAnalytics(c, db)
	})
	// Serve uploaded files
	router.Use(static.Serve("/", static.LocalFile("./uploads", true)))

	// File upload endpoints
	router.POST("/upload", handlers.UploadFile)
	router.GET("/files", handlers.GetUploadedFiles)
	//company router
	router.POST("/create-company", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.CreateCompany(c, db)
	})
	router.GET("/companies/:companyId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.GetCompanyByID(c, db)
	})
	router.PUT("/companies/:companyId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.UpdateCompany(c, db)
	})
	router.DELETE("/companies/:companyId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.DeleteCompany(c, db)
	})
	// Role routes
	router.POST("/roles", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.CreateRole(c, db)
	})

	router.PUT("/roles/:roleId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.EditRole(c, db)
	})

	router.DELETE("/roles/:roleId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.DeleteRole(c, db)
	})

	router.GET("/roles", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.GetAllRoles(c, db)
	})

	router.GET("/roles/:roleId", handlers.AuthMiddleware(), func(c *gin.Context) {
		handlers.GetRoleByID(c, db)
	})
	// Debugging route
	router.GET("/debug/routes", func(c *gin.Context) {
		fmt.Println("yes")
		c.JSON(http.StatusOK, router.Routes())
	})

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
