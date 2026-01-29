package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMe 現在のユーザー情報を取得（保護されたエンドポイント）
func GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	email, _ := c.Get("email")
	username, _ := c.Get("username")

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       userID,
			"email":    email,
			"username": username,
		},
	})
}

// GetProfile ユーザープロファイルを取得（保護されたエンドポイント）
func GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	userMutex.RLock()
	user, exists := users[userID.(uint)]
	userMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}
