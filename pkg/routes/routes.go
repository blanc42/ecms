package routes

import (
	"github.com/blanc42/ecms/pkg/handlers"
	"github.com/blanc42/ecms/pkg/initializers"
	"github.com/blanc42/ecms/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(re *gin.Engine) {

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "https://dhukan.vercel.app"} // Add your frontend URL
	config.AllowCredentials = true
	re.Use(cors.New(config))
	re.RedirectTrailingSlash = false

	r := re.Group("/api/")

	AdminHandler := handlers.NewAdminHandler(initializers.DB)
	storeHandler := handlers.NewStoreHandler(initializers.DB)
	categoryHandler := handlers.NewCategoryHandler(initializers.DB)

	r.POST("/signup", AdminHandler.Signup)
	r.POST("/login", AdminHandler.Login)

	adminOnly := r.Group("/")
	adminOnly.Use(middleware.AdminAuthMiddleware())
	adminOnly.GET("/admin", AdminHandler.GetAdmin)

	storeGroup := r.Group("/stores")
	storeGroup.Use(middleware.AdminAuthMiddleware())
	{
		adminOnly.POST("/stores", storeHandler.CreateStore)
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
