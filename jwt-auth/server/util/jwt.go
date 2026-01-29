package util

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/makoto-developer/golang_examples/jwt-auth/server/model"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GetJWTSecret JWT秘密鍵を取得
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-this"
	}
	return []byte(secret)
}

// GetAccessExpires アクセストークン有効期限を取得（秒）
func GetAccessExpires() int64 {
	expires := os.Getenv("JWT_ACCESS_EXPIRES")
	if expires == "" {
		return 3600 // デフォルト: 1時間
	}
	exp, err := strconv.ParseInt(expires, 10, 64)
	if err != nil {
		return 3600
	}
	return exp
}

// GetRefreshExpires リフレッシュトークン有効期限を取得（秒）
func GetRefreshExpires() int64 {
	expires := os.Getenv("JWT_REFRESH_EXPIRES")
	if expires == "" {
		return 604800 // デフォルト: 7日間
	}
	exp, err := strconv.ParseInt(expires, 10, 64)
	if err != nil {
		return 604800
	}
	return exp
}

// GenerateAccessToken アクセストークンを生成
func GenerateAccessToken(claims model.Claims) (string, error) {
	expiresIn := GetAccessExpires()
	jwtClaims := JWTClaims{
		UserID:   claims.UserID,
		Email:    claims.Email,
		Username: claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(GetJWTSecret())
}

// GenerateRefreshToken リフレッシュトークンを生成
func GenerateRefreshToken(claims model.Claims) (string, error) {
	expiresIn := GetRefreshExpires()
	jwtClaims := JWTClaims{
		UserID:   claims.UserID,
		Email:    claims.Email,
		Username: claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(GetJWTSecret())
}

// ValidateToken トークンを検証
func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return GetJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// 有効期限チェック
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
