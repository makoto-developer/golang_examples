package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/makoto-developer/golang_examples/jwt-auth/server/handler"
	"github.com/makoto-developer/golang_examples/jwt-auth/server/middleware"
)

func main() {
	// .envファイル読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Ginモード設定
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// CORS設定（開発用）
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// ヘルスチェック
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1
	v1 := router.Group("/api")
	{
		// 認証エンドポイント（公開）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
			auth.POST("/refresh", handler.Refresh)
		}

		// ユーザーエンドポイント（保護）
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/me", handler.GetMe)
			users.GET("/profile", handler.GetProfile)
		}
	}

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "22500"
	}

	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
