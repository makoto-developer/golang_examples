package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/makoto-developer/golang_examples/gorm/gorm/model"
	"github.com/makoto-developer/golang_examples/gorm/gorm/repository"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config はアプリケーション設定を保持
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

// DatabaseConfig はデータベース設定
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ServerConfig はサーバー設定
type ServerConfig struct {
	Port string
}

var (
	db        *gorm.DB
	orderRepo repository.OrderRepository
	config    *Config
)

// initConfig は設定ファイルをロード
func initConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // 環境変数を自動的に読み込む

	// デフォルト値を設定
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("SERVER_PORT", "8080")

	// 設定ファイルをロード（存在しなくてもエラーにしない）
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Warning: Error reading config file: %v", err)
		} else {
			log.Println("Warning: .env file not found, using environment variables")
		}
	} else {
		log.Printf("Loaded config from: %s", viper.ConfigFileUsed())
	}

	// 設定を構造体にマッピング（大文字のキーを使用）
	config = &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
	}

	// パスワードが設定されているか確認
	if config.Database.Password == "" {
		log.Fatalf("Database password is not set. Please set DB_PASSWORD in .env file or environment variable.")
	}

	log.Printf("Config loaded: Host=%s, Port=%s, User=%s, DBName=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.DBName,
	)
}

// initDB はデータベース接続を初期化
func initDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Database.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.Port,
		config.Database.SSLMode,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
}

// initRepository はRepositoryを初期化
func initRepository() {
	orderRepo = repository.NewOrderRepository(db)
	log.Println("Repository initialized successfully")
}

func init() {
	initConfig()
	initDB()
	initRepository()
}

func main() {
	r := gin.Default()

	// ヘルスチェック
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "fine!",
		})
	})

	// 注文一覧取得
	r.GET("/orders", getOrders)

	// 注文詳細取得
	r.GET("/orders/:id", getOrder)

	// ユーザーIDで注文検索
	r.GET("/users/:user_id/orders", getOrdersByUserID)

	// 新規注文作成
	r.POST("/orders", createOrder)

	// 注文更新
	r.PUT("/orders/:id", updateOrder)

	// 注文削除
	r.DELETE("/orders/:id", deleteOrder)

	r.Run(":" + config.Server.Port)
}

// getOrders は全注文を取得
func getOrders(c *gin.Context) {
	ids := c.Query("ids")

	if ids != "" {
		var orderIDs []uint64
		idStrings := strings.Split(ids, ",")

		for _, idStr := range idStrings {
			id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Invalid ID format: %s", idStr),
				})
				return
			}
			orderIDs = append(orderIDs, id)
		}

		orders, err := orderRepo.ListByOrderID(orderIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, orders)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Use ?ids=1,2,3 to get specific orders",
	})
}

// getOrder は指定IDの注文を取得
func getOrder(c *gin.Context) {
	id := c.Param("id")
	orderID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	order := orderRepo.Get(orderID)
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// getOrdersByUserID はユーザーIDで注文を検索
func getOrdersByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	orders, err := orderRepo.ListByUserID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// createOrder は新規注文を作成
func createOrder(c *gin.Context) {
	var order model.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	if err := orderRepo.Create(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// updateOrder は注文を更新
func updateOrder(c *gin.Context) {
	id := c.Param("id")
	orderID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	order := orderRepo.Get(orderID)
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	if err := orderRepo.Update(*order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// deleteOrder は注文を削除
func deleteOrder(c *gin.Context) {
	id := c.Param("id")
	orderID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	if err := orderRepo.Delete(orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order deleted successfully",
	})
}
