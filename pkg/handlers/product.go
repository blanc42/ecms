package handlers

import (
	"net/http"
	"strconv"

	"github.com/blanc42/ecms/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

type ProductItemInput struct {
	SKU             string  `json:"sku" binding:"required"`
	Quantity        int     `json:"quantity" binding:"required,gte=0"`
	Price           float64 `json:"price" binding:"required,gt=0"`
	DiscountedPrice float64 `json:"discounted_price,omitempty"`
}

type CreateProductInput struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	Rating      float32            `json:"rating"`
	IsFeatured  bool               `json:"is_featured"`
	IsArchived  bool               `json:"is_archived"`
	HasVariants bool               `json:"has_variants"`
	CategoryID  uint               `json:"category_id" binding:"required"`
	StoreID     uint               `json:"store_id" binding:"required"`
	Items       []ProductItemInput `json:"items"`
}

type UpdateProductItemInput struct {
	ID              *uint   `json:"id"`
	SKU             string  `json:"sku" binding:"required"`
	Quantity        int     `json:"quantity" binding:"required,gte=0"`
	Price           float64 `json:"price" binding:"required,gt=0"`
	DiscountedPrice float64 `json:"discounted_price,omitempty"`
}

type UpdateProductInput struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Rating      float32                  `json:"rating"`
	IsFeatured  bool                     `json:"is_featured"`
	IsArchived  bool                     `json:"is_archived"`
	HasVariants bool                     `json:"has_variants"`
	CategoryID  uint                     `json:"category_id"`
	Items       []UpdateProductItemInput `json:"items"`
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the category exists and belongs to the current admin's store
	var category models.Category
	if err := h.DB.First(&category, input.CategoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	adminID, _ := c.Get("admin_id")
	var store models.Store
	if err := h.DB.First(&store, input.StoreID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch store"})
		return
	}

	if store.AdminID != adminID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to create products for this store"})
		return
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Rating:      input.Rating,
		IsFeatured:  input.IsFeatured,
		IsArchived:  input.IsArchived,
		HasVariants: input.HasVariants,
		CategoryID:  input.CategoryID,
		StoreID:     input.StoreID,
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		for _, itemInput := range input.Items {
			item := models.ProductItem{
				ProductID:       product.ID,
				SKU:             itemInput.SKU,
				Quantity:        itemInput.Quantity,
				Price:           itemInput.Price,
				DiscountedPrice: itemInput.DiscountedPrice,
			}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": product, "error": nil})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("product_id")
	var product models.Product

	if err := h.DB.Preload("Items").First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product, "error": nil})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID := c.Param("product_id")
	var product models.Product

	if err := h.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Update product fields if provided
		if input.Name != "" {
			product.Name = input.Name
		}
		if input.Description != "" {
			product.Description = input.Description
		}
		product.Rating = input.Rating
		product.IsFeatured = input.IsFeatured
		product.IsArchived = input.IsArchived
		product.HasVariants = input.HasVariants
		if input.CategoryID != 0 {
			product.CategoryID = input.CategoryID
		}

		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		// Handle product items
		for _, itemInput := range input.Items {
			if itemInput.ID != nil {
				// Update existing item
				var existingItem models.ProductItem
				if err := tx.First(&existingItem, *itemInput.ID).Error; err != nil {
					return err
				}
				existingItem.SKU = itemInput.SKU
				existingItem.Quantity = itemInput.Quantity
				existingItem.Price = itemInput.Price
				existingItem.DiscountedPrice = itemInput.DiscountedPrice
				if err := tx.Save(&existingItem).Error; err != nil {
					return err
				}
			} else {
				// Create new item
				newItem := models.ProductItem{
					ProductID:       product.ID,
					SKU:             itemInput.SKU,
					Quantity:        itemInput.Quantity,
					Price:           itemInput.Price,
					DiscountedPrice: itemInput.DiscountedPrice,
				}
				if err := tx.Create(&newItem).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Fetch the updated product with its items
	if err := h.DB.Preload("Items").First(&product, productID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product, "error": nil})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID := c.Param("product_id")
	var product models.Product

	if err := h.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Delete associated product items
		if err := tx.Where("product_id = ?", product.ID).Delete(&models.ProductItem{}).Error; err != nil {
			return err
		}

		// Delete the product
		if err := tx.Delete(&product).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	storeID := c.Param("store_id")
	var products []models.Product

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	if err := h.DB.Where("store_id = ?", storeID).Preload("Items").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products, "error": nil})
}


func (h *ProductHandler) GetFilters(c *gin.Context){
		
}