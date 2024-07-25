package routes

import (
	"github.com/blanc42/ecms/pkg/handlers"
	"github.com/blanc42/ecms/pkg/initializers"
	"github.com/blanc42/ecms/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	AdminHandler := handlers.NewAdminHandler(initializers.DB)
	storeHandler := handlers.NewStoreHandler(initializers.DB)
	categoryHandler := handlers.NewCategoryHandler(initializers.DB)

	r.POST("/signup", AdminHandler.Signup)
	r.POST("/login", AdminHandler.Login)

	adminOnly := r.Group("/")
	adminOnly.Use(middleware.AdminAuthMiddleware())

	storeGroup := r.Group("/stores")
	storeGroup.Use(middleware.AdminAuthMiddleware())
	{
		storeGroup.POST("/", storeHandler.CreateStore)
		storeGroup.GET("/:store_id", storeHandler.GetStore)
		storeGroup.PUT("/:store_id", storeHandler.UpdateStore)
		storeGroup.DELETE("/:store_id", storeHandler.DeleteStore)
		storeGroup.GET("/", storeHandler.ListStores)
	}

	{
		storeGroup.POST("/:store_id/categories", categoryHandler.CreateCategory)
		storeGroup.GET("/:store_id/categories", categoryHandler.GetAllCategories)
		storeGroup.GET("/:store_id/categories/:category_id", categoryHandler.GetCategory)
		storeGroup.PUT("/:store_id/categories/:category_id", categoryHandler.UpdateCategory)
		storeGroup.DELETE("/:store_id/categories/:category_id", categoryHandler.DeleteCategory)
	}

	variantHandler := handlers.NewVariantHandler(initializers.DB)

	storeGroup.POST("/:store_id/categories/:category_id/variants", variantHandler.CreateVariant)
	storeGroup.GET("/:store_id/categories/:category_id/variants", variantHandler.ListVariants)
	storeGroup.GET("/:store_id/categories/:category_id/variants/:variant_id", variantHandler.GetVariant)
	storeGroup.PUT("/:store_id/categories/:category_id/variants/:variant_id", variantHandler.UpdateVariant)
	storeGroup.DELETE("/:store_id/categories/:category_id/variants/:variant_id", variantHandler.DeleteVariant)

	productHandler := handlers.NewProductHandler(initializers.DB)

	storeGroup.POST("/:store_id/products", productHandler.CreateProduct)
	storeGroup.GET("/:store_id/products", productHandler.ListProducts)
	storeGroup.GET("/:store_id/products/:product_id", productHandler.GetProduct)
	storeGroup.PUT("/:store_id/products/:product_id", productHandler.UpdateProduct)
	storeGroup.DELETE("/:store_id/products/:product_id", productHandler.DeleteProduct)

}
