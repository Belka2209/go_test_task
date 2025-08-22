package main

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
)

func main() {
	// Настройка использования максимального количества CPU
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Настройка роутера
	router := gin.Default()

	// Маршруты аутентификации
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", handleLogin)
		authGroup.POST("/register", handleRegister)
	}

	// Защищенные маршруты
	protected := router.Group("/api")
	protected.Use(authMiddleware)
	{
		protected.GET("/userid", handleGetUserID)
		protected.POST("/analytics", handleItemAnalytics)
	}

	// Корневой маршрут
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Microservice is running!"})
	})

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": "test-token", "user_id": "test-user"})
}

func handleRegister(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func authMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		c.Abort()
		return
	}

	// Простая проверка токена
	if token != "Bearer test-token" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Next()
}

func handleGetUserID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"user_id": "test-user"})
}

func handleItemAnalytics(c *gin.Context) {
	var request struct {
		Items []struct {
			ID         string    `json:"id"`
			Sales      []float64 `json:"sales"`
			StockLevel []float64 `json:"stock_level"`
			Price      float64   `json:"price"`
		} `json:"items"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Простая аналитика
	results := make([]gin.H, len(request.Items))
	for i, item := range request.Items {
		score := 0.0
		if len(item.Sales) > 0 {
			for _, sale := range item.Sales {
				score += sale
			}
			score /= float64(len(item.Sales))
		}
		score *= item.Price

		results[i] = gin.H{
			"item_id": item.ID,
			"score":   score,
			"status":  "success",
		}
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
