package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	rdb    *redis.Client
	log    *logrus.Logger
	config *Config
)

// Config holds application configuration
type Config struct {
	RegistryURL  string
	RegistryPort string
	APIPort      string
	DBHost       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPassword   string
	RedisHost    string
	RedisPort    string
	JWTSecret    string
}

// User represents a registry user
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	Email    string
	Role     string
	Active   bool
}

// Image represents a Docker image
type Image struct {
	gorm.Model
	Name        string
	Description string
	Tags        []Tag
	Owner       string
	Public      bool
}

// Tag represents an image tag
type Tag struct {
	gorm.Model
	ImageID     uint
	Name        string
	Digest      string
	Size        int64
	CreatedAt   time.Time
	LastPulled  time.Time
	PullCount   int64
}

func init() {
	// Initialize logger
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	
	// Load configuration
	config = loadConfig()

	// Initialize database
	initDB()

	// Initialize Redis
	initRedis()
}

func loadConfig() *Config {
	return &Config{
		RegistryURL:  getEnv("REGISTRY_URL", "localhost"),
		RegistryPort: getEnv("REGISTRY_PORT", "5000"),
		APIPort:      getEnv("API_PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBName:       getEnv("DB_NAME", "registry"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "postgres"),
		RedisHost:    getEnv("REDIS_HOST", "localhost"),
		RedisPort:    getEnv("REDIS_PORT", "6379"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func initDB() {
	dsn := "host=" + config.DBHost + " user=" + config.DBUser +
		" password=" + config.DBPassword + " dbname=" + config.DBName +
		" port=" + config.DBPort + " sslmode=disable TimeZone=UTC"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&User{}, &Image{}, &Tag{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Info("Database initialized successfully")
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test Redis connection
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Info("Redis initialized successfully")
}

func main() {
	// Create Gin router
	router := gin.Default()

	// Configure trusted proxies
	router.SetTrustedProxies([]string{"127.0.0.1", "172.19.0.0/16"})

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// API version group
	v1 := router.Group("/api/v1")
	{
		// Registry routes
		registry := v1.Group("/registry")
		{
			registry.GET("/health", healthCheck)
			registry.GET("/info", getRegistryInfo)
			registry.GET("/metrics", getRegistryMetrics)
		}

		// Image routes
		images := v1.Group("/images")
		{
			images.GET("", listImages)
			images.GET("/:name", getImage)
			images.GET("/:name/tags", listImageTags)
			images.DELETE("/:name", deleteImage)
			images.DELETE("/:name/tags/:tag", deleteImageTag)
		}

		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", login)
			auth.POST("/token", getToken)
			auth.GET("/verify", verifyToken)
		}

		// User routes
		users := v1.Group("/users")
		{
			users.GET("", listUsers)
			users.POST("", createUser)
			users.PUT("/:username", updateUser)
			users.DELETE("/:username", deleteUser)
		}
	}

	// Start server
	port := ":" + config.APIPort
	log.Infof("Starting server on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Handler function placeholders
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "UP",
		"timestamp": time.Now().UTC(),
	})
}

func getRegistryInfo(c *gin.Context) {
	// Get registry information from the Docker Registry API
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/v2/", config.RegistryURL, config.RegistryPort))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to registry"})
		return
	}
	defer resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{
		"version": "2.0",
		"url": fmt.Sprintf("%s:%s", config.RegistryURL, config.RegistryPort),
		"status": "running",
	})
}

func getRegistryMetrics(c *gin.Context) {
	// Get metrics from Redis
	ctx := context.Background()
	imageCount, err := rdb.Get(ctx, "registry:image_count").Int()
	if err != nil {
		imageCount = 0
	}

	tagCount, err := rdb.Get(ctx, "registry:tag_count").Int()
	if err != nil {
		tagCount = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"images": imageCount,
		"tags": tagCount,
		"timestamp": time.Now().UTC(),
	})
}

func listImages(c *gin.Context) {
	var images []Image
	result := db.Find(&images)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
		return
	}

	c.JSON(http.StatusOK, images)
}

func getImage(c *gin.Context) {
	name := c.Param("name")
	var image Image
	
	result := db.Where("name = ?", name).First(&image)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.JSON(http.StatusOK, image)
}

func listImageTags(c *gin.Context) {
	name := c.Param("name")
	var image Image
	var tags []Tag

	if err := db.Where("name = ?", name).First(&image).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	if err := db.Where("image_id = ?", image.ID).Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func deleteImage(c *gin.Context) {
	name := c.Param("name")
	var image Image

	if err := db.Where("name = ?", name).First(&image).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Delete associated tags first
	if err := db.Where("image_id = ?", image.ID).Delete(&Tag{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image tags"})
		return
	}

	// Delete the image
	if err := db.Delete(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}

func deleteImageTag(c *gin.Context) {
	name := c.Param("name")
	tagName := c.Param("tag")
	var image Image
	var tag Tag

	if err := db.Where("name = ?", name).First(&image).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	if err := db.Where("image_id = ? AND name = ?", image.ID, tagName).First(&tag).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	if err := db.Delete(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// TODO: Implement proper password hashing
	if req.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": user,
	})
}

func getToken(c *gin.Context) {
	// This endpoint is for refreshing tokens
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Generate new token
	newToken := jwt.New(jwt.SigningMethodHS256)
	newClaims := newToken.Claims.(jwt.MapClaims)
	newClaims["username"] = user.Username
	newClaims["role"] = user.Role
	newClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	newTokenString, err := newToken.SignedString([]byte(config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newTokenString})
}

func verifyToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"claims": claims,
	})
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
}

func listUsers(c *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if username already exists
	var existingUser User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Create new user
	user := User{
		Username: req.Username,
		Password: req.Password, // TODO: Hash password
		Email:    req.Email,
		Role:     req.Role,
		Active:   true,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
	username := c.Param("username")
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Email = req.Email
	user.Role = req.Role
	if req.Password != "" {
		user.Password = req.Password // TODO: Hash password
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	username := c.Param("username")
	var user User

	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
