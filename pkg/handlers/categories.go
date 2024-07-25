/*
- check if the parent category exists when creating a category and other checks on it
  - this will be taken care by the fk constraint
  - but it's not working when we are creating the category with id 1 and parent_category_id = 1 (self referencing)

- check if the store exists
- check if the admin is the owner of the store
*/
package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/blanc42/ecms/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	DB *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

type CreateCategoryInput struct {
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description"`
	StoreID          uint   `json:"store_id" binding:"required"`
	ParentCategoryID *uint  `json:"parent_category_id,omitempty"`
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the store exists and belongs to the current admin
	var store models.Store
	if err := h.DB.First(&store, input.StoreID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	adminID, _ := c.Get("admin_id")
	if store.AdminID != adminID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to create categories for this store"})
		return
	}

	category := models.Category{
		Name:             input.Name,
		Description:      input.Description,
		StoreID:          input.StoreID,
		ParentCategoryID: input.ParentCategoryID,
	}

	if err := h.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": category, "error": nil})
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	categoryID := c.Param("category_id")
	var category models.Category

	if err := h.DB.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category, "error": nil})
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	categoryID := c.Param("category_id")
	var category models.Category

	if err := h.DB.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var input CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Model(&category).Updates(models.Category{
		Name:             input.Name,
		Description:      input.Description,
		ParentCategoryID: input.ParentCategoryID,
	})

	c.JSON(http.StatusOK, gin.H{"data": category, "error": nil})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryID := c.Param("category_id")
	var category models.Category

	if err := h.DB.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	h.DB.Delete(&category)
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	storeID, err := strconv.ParseUint(c.Param("store_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid store ID"})
		fmt.Println(err)
		return
	}

	level, _ := strconv.Atoi(c.DefaultQuery("level", "-1"))

	var categories []models.Category

	// Use a CTE to fetch the entire category tree in a single query
	query := h.DB.Raw(`
		WITH RECURSIVE category_tree AS (
			SELECT *, 0 AS level
			FROM categories
			WHERE store_id = ? AND parent_category_id IS NULL
			
			UNION ALL
			
			SELECT c.*, ct.level + 1
			FROM categories c
			JOIN category_tree ct ON c.parent_category_id = ct.id
			WHERE c.store_id = ?
		)
		SELECT *
		FROM category_tree
		WHERE ? = -1 OR level <= ?
		ORDER BY id
	`, storeID, storeID, level, level)

	if err := query.Scan(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	// Convert flat structure to tree
	categoryMap := make(map[uint]*models.Category)
	var result []*models.Category

	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
		if categories[i].ParentCategoryID == nil {
			result = append(result, &categories[i])
		} else {
			parent := categoryMap[*categories[i].ParentCategoryID]
			parent.Subcategories = append(parent.Subcategories, &categories[i])
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": result, "error": nil})
}
