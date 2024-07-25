package handlers

import (
	"net/http"
	"strconv"

	"github.com/blanc42/ecms/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StoreHandler struct {
	DB *gorm.DB
}

func NewStoreHandler(db *gorm.DB) *StoreHandler {
	return &StoreHandler{DB: db}
}

type CreateStoreInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *StoreHandler) CreateStore(c *gin.Context) {
	var input CreateStoreInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, _ := c.Get("admin_id")
	store := models.Store{
		Name:        input.Name,
		Description: input.Description,
		AdminID:     adminID.(uint),
	}

	if err := h.DB.Create(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create store"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": store, "error": nil})
}
func (h *StoreHandler) GetStore(c *gin.Context) {
	storeID := c.Param("store_id")
	var store models.Store

	if err := h.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	adminID, _ := c.Get("admin_id")
	if store.AdminID != adminID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this store"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": store, "error": nil})
}

func (h *StoreHandler) UpdateStore(c *gin.Context) {
	storeID := c.Param("store_id")
	var store models.Store

	if err := h.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	adminID, _ := c.Get("admin_id")
	if store.AdminID != adminID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this store"})
		return
	}

	var input CreateStoreInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Model(&store).Updates(models.Store{Name: input.Name, Description: input.Description})
	c.JSON(http.StatusOK, gin.H{"data": store, "error": nil})
}

func (h *StoreHandler) DeleteStore(c *gin.Context) {
	storeID := c.Param("store_id")
	var store models.Store

	if err := h.DB.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	adminID, _ := c.Get("admin_id")
	if store.AdminID != adminID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this store"})
		return
	}

	h.DB.Delete(&store)
	c.JSON(http.StatusOK, gin.H{"message": "Store deleted successfully"})
}

func (h *StoreHandler) ListStores(c *gin.Context) {
	var stores []models.Store
	adminID, _ := c.Get("admin_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	h.DB.Where("admin_id = ?", adminID).Offset(offset).Limit(pageSize).Find(&stores)

	c.JSON(http.StatusOK, gin.H{"data": stores, "error": nil})
}
