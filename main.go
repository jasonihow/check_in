package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User model with real name added
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
}

// CheckIn model to store user check-in data
type CheckIn struct {
	gorm.Model
	Username    string    `json:"username"`      // Username used for check-in
	CheckInTime time.Time `json:"check_in_time"` // Time of check-in
}

var userDB *gorm.DB

func main() {
	var err error
	userDB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("cannot connect to database")
	}

	// Migrate the User and CheckIn models
	userDB.AutoMigrate(&User{}, &CheckIn{})

	r := gin.Default()

	r.POST("/register", register)
	r.POST("/login", login)
	r.GET("/users", getUsers)
	r.POST("/check_in/:username", checkIn)
	r.GET("/check_ins", getCheckIns)

	r.Run(":8080")
}

// Register a new user
func register(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": ErrMissingParameter, "error": GetErrorMessage(ErrMissingParameter)})
		return
	}

	// Check if the username already exists
	var existingUser User
	if err := userDB.Where("username = ?", json.Username).First(&existingUser).Error; err == nil {
		// If err is nil, that means a user with this username already exists
		c.JSON(http.StatusBadRequest, gin.H{"error_code": ErrUsingSameUsername, "error": "Username already exists"})
		return
	}

	// Create a new user record
	if err := userDB.Create(&json).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": ErrUserCreationFailed, "error": GetErrorMessage(ErrUserCreationFailed)})
		return
	}

	// Ensure we return the name in the message
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "User added: " + json.Username,
	})
}

// Get all users
func getUsers(c *gin.Context) {
	var users []User

	// Fetch all users from the database
	if err := userDB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": ErrFetchingUsersFailed, "error": GetErrorMessage(ErrFetchingUsersFailed)})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Login function
func login(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": ErrLoginMissingParameter, "error": GetErrorMessage(ErrLoginMissingParameter)})
		return
	}

	var user User
	if err := userDB.Where("username = ? AND password = ?", json.Username, json.Password).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error_code": ErrInvalidCredentials, "error": GetErrorMessage(ErrInvalidCredentials)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": ErrLoginDatabaseFailed, "error": GetErrorMessage(ErrLoginDatabaseFailed)})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Record user check-in time based on username from URL parameter
func checkIn(c *gin.Context) {
	username := c.Param("username")

	// Verify user exists
	var user User
	if err := userDB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error_code": ErrUserNotFound, "error": GetErrorMessage(ErrUserNotFound)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error_code": ErrLoginDatabaseFailed, "error": GetErrorMessage(ErrLoginDatabaseFailed)})
		}
		return
	}

	// Create a new check-in record
	checkIn := CheckIn{
		Username:    username,
		CheckInTime: time.Now(),
	}
	if err := userDB.Create(&checkIn).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": ErrCheckInFailed, "error": GetErrorMessage(ErrCheckInFailed)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Check-in successful",
		"check_in_id":   checkIn.ID,
		"username":      checkIn.Username,
		"check_in_time": checkIn.CheckInTime.Format(time.RFC3339),
	})
}

// Get all check-in records
func getCheckIns(c *gin.Context) {
	var checkIns []CheckIn

	// Fetch all check-ins from the database
	if err := userDB.Find(&checkIns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": ErrFetchingCheckInsFailed, "error": GetErrorMessage(ErrFetchingCheckInsFailed)})
		return
	}

	c.JSON(http.StatusOK, checkIns)
}
