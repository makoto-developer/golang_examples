package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/makoto-developer/golang_examples/jwt-auth/server/model"
	"github.com/makoto-developer/golang_examples/jwt-auth/server/util"
)

// インメモリデータストア（本番環境ではデータベースを使用）
var (
	users      = make(map[uint]*model.User)
	usersByEmail = make(map[string]*model.User)
	userID     uint = 1
	userMutex  sync.RWMutex
)

// Register ユーザー登録
func Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userMutex.Lock()
	defer userMutex.Unlock()

	// メールアドレス重複チェック
	if _, exists := usersByEmail[req.Email]; exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already registered",
		})
		return
	}

	// パスワードハッシュ化
	hashedPassword, err := model.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// ユーザー作成
	user := &model.User{
		ID:        userID,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	users[userID] = user
	usersByEmail[req.Email] = user
	userID++

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login ログイン
func Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userMutex.RLock()
	user, exists := usersByEmail[req.Email]
	userMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// パスワード確認
	if !model.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// トークン生成
	claims := model.Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	accessToken, err := util.GenerateAccessToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	refreshToken, err := util.GenerateRefreshToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate refresh token",
		})
		return
	}

	response := model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    util.GetAccessExpires(),
		TokenType:    "Bearer",
	}

	c.JSON(http.StatusOK, response)
}

// Refresh トークンリフレッシュ
func Refresh(c *gin.Context) {
	var req model.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// リフレッシュトークン検証
	claims, err := util.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired refresh token",
		})
		return
	}

	// 新しいアクセストークン生成
	newClaims := model.Claims{
		UserID:   claims.UserID,
		Email:    claims.Email,
		Username: claims.Username,
	}

	accessToken, err := util.GenerateAccessToken(newClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	response := model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken, // リフレッシュトークンは再利用
		ExpiresIn:    util.GetAccessExpires(),
		TokenType:    "Bearer",
	}

	c.JSON(http.StatusOK, response)
}
