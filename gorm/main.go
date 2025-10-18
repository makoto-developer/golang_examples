package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	gormConfig "github.com/makoto-developer/golang_examples/gorm/config"
	"github.com/makoto-developer/golang_examples/gorm/gorm/handler"
	"github.com/makoto-developer/golang_examples/gorm/gorm/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db           *gorm.DB
	orderRepo    repository.OrderRepository
	orderHandler *handler.OrderHandler
	config       *gormConfig.Config
)

func initDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
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

func initRepository() {
	orderRepo = repository.NewOrderRepository(db)
	log.Println("Repository initialized successfully")
}

func initHandler() {
	orderHandler = handler.NewOrderHandler(orderRepo)
	log.Println("Handler initialized successfully")
}

func init() {
	config = gormConfig.LoadConfig()
	initDB()
	initRepository()
	initHandler()
}

func main() {
	r := gin.Default()
	setupRoutes(r)
	err := r.Run(":" + config.Server.Port)
	if err != nil {
		_ = fmt.Errorf("failed to start server: %v", err)
		return
	}
}

func setupRoutes(r *gin.Engine) {
	r.GET("/ping", healthCheck)

	orders := r.Group("/orders")
	{
		orders.GET("", orderHandler.GetOrders)
		orders.GET("/:id", orderHandler.GetOrder)
		orders.POST("", orderHandler.CreateOrder)
		orders.PUT("/:id", orderHandler.UpdateOrder)
		orders.DELETE("/:id", orderHandler.DeleteOrder)
	}

	users := r.Group("/users")
	{
		users.GET("/:user_id/orders", orderHandler.GetOrdersByUserID)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "fine!"})
}
