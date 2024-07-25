/*
unique indie for variant_id, option_value
for example variant_id = 1, option_value = red, blue, green
we can add `red` option again and again with a different option id
*/

package handlers

import (
	"net/http"
	"strconv"

	"github.com/blanc42/ecms/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VariantHandler struct {
	DB *gorm.DB
}

func NewVariantHandler(db *gorm.DB) *VariantHandler {
	return &VariantHandler{DB: db}
}

type VariantOptionInput struct {
	Value       string `json:"value" binding:"required"`
	Description string `json:"description"`
	Weight      int    `json:"weight"`
}

type CreateVariantInput struct {
	Name        string               `json:"name" binding:"required"`
	Description string               `json:"description"`
	Weight      int                  `json:"weight" binding:"required,gt=0"`
	CategoryID  uint                 `json:"category_id" binding:"required"`
	Options     []VariantOptionInput `json:"options"`
}

type UpdateVariantOptionInput struct {
	ID          *uint  `json:"id"`
	Value       string `json:"value" binding:"required"`
	Description string `json:"description"`
	Weight      int    `json:"weight"`
}

type UpdateVariantInput struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Weight      int                        `json:"weight"`
	CategoryID  uint                       `json:"category_id"`
	Options     []UpdateVariantOptionInput `json:"options"`
}

func (h *VariantHandler) CreateVariant(c *gin.Context) {
	var input CreateVariantInput
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
	if err := h.DB.First(&store, category.StoreID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch store"})
		return
	}

	if store.AdminID != adminID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to create variants for this category"})
		return
	}

	variant := models.Variant{
		Name:        input.Name,
		Description: input.Description,
		Weight:      input.Weight,
		CategoryID:  input.CategoryID,
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&variant).Error; err != nil {
			return err
		}

		for _, optionInput := range input.Options {
			option := models.VariantOption{
				Value:       optionInput.Value,
				Description: optionInput.Description,
				Weight:      optionInput.Weight,
				VariantID:   variant.ID,
			}
			if err := tx.Create(&option).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create variant"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": variant, "error": nil})
}

func (h *VariantHandler) GetVariant(c *gin.Context) {
	variantID := c.Param("variant_id")
	var variant models.Variant

	if err := h.DB.Preload("Options").First(&variant, variantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": variant, "error": nil})
}

func (h *VariantHandler) UpdateVariant(c *gin.Context) {
	variantID := c.Param("variant_id")
	var variant models.Variant

	if err := h.DB.First(&variant, variantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
		return
	}

	var input UpdateVariantInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Update variant fields if provided
		if input.Name != "" {
			variant.Name = input.Name
		}
		if input.Description != "" {
			variant.Description = input.Description
		}
		if input.Weight != 0 {
			variant.Weight = input.Weight
		}
		if input.CategoryID != 0 {
			variant.CategoryID = input.CategoryID
		}

		if err := tx.Save(&variant).Error; err != nil {
			return err
		}

		// Handle options
		for _, optionInput := range input.Options {
			if optionInput.ID != nil {
				// Update existing option
				var existingOption models.VariantOption
				if err := tx.First(&existingOption, *optionInput.ID).Error; err != nil {
					return err
				}
				existingOption.Value = optionInput.Value
				existingOption.Description = optionInput.Description
				existingOption.Weight = optionInput.Weight
				if err := tx.Save(&existingOption).Error; err != nil {
					return err
				}
			} else {
				// Create new option
				newOption := models.VariantOption{
					Value:       optionInput.Value,
					Description: optionInput.Description,
					Weight:      optionInput.Weight,
					VariantID:   variant.ID,
				}
				if err := tx.Create(&newOption).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update variant"})
		return
	}

	// Fetch the updated variant with its options
	if err := h.DB.Preload("Options").First(&variant, variantID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated variant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": variant, "error": nil})
}

func (h *VariantHandler) DeleteVariant(c *gin.Context) {
	variantID := c.Param("variant_id")
	var variant models.Variant

	if err := h.DB.First(&variant, variantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Delete associated options
		if err := tx.Where("variant_id = ?", variant.ID).Delete(&models.VariantOption{}).Error; err != nil {
			return err
		}

		// Delete the variant
		if err := tx.Delete(&variant).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete variant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Variant deleted successfully"})
}

func (h *VariantHandler) ListVariants(c *gin.Context) {
	categoryID := c.Param("category_id")
	var variants []models.Variant

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	if err := h.DB.Where("category_id = ?", categoryID).Preload("Options").Offset(offset).Limit(pageSize).Find(&variants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch variants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": variants, "error": nil})
}
