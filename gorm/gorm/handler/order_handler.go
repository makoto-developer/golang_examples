package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/makoto-developer/golang_examples/gorm/gorm/model"
	"github.com/makoto-developer/golang_examples/gorm/gorm/repository"
)

type OrderHandler struct {
	repo repository.OrderRepository
}

func NewOrderHandler(repo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
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

		orders, err := h.repo.ListByOrderID(orderIDs)
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

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	orderID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	order := h.repo.Get(orderID)
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrdersByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	orders, err := h.repo.ListByUserID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order model.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	if err := h.validateOrder(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.repo.Create(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	orderID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	order := h.repo.Get(orderID)
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

	if err := h.validateOrder(order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.repo.Update(*order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	orderID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	if err := h.repo.Delete(orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order deleted successfully",
	})
}

func (h *OrderHandler) validateOrder(order *model.Order) error {
	if order.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}

	if order.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if order.AmountWithoutTax < 0 {
		return fmt.Errorf("amount_without_tax cannot be negative")
	}

	if order.Tax < 0 {
		return fmt.Errorf("tax cannot be negative")
	}

	return nil
}
