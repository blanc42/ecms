package initializers

import (
	"log"
	"os"

	"github.com/blanc42/ecms/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Migrate the schema
	err = DB.AutoMigrate(
		&models.Admin{},
		&models.Store{},
		&models.Category{},
		&models.Product{},
		&models.Variant{},
		&models.VariantOption{},
		&models.ProductItem{},
		&models.ProductImage{},
		// &models.Customer{},
		// &models.Order{},
		// &models.OrderItem{},
		// &models.Cart{},
		// &models.CartItem{},
		// &models.Address{},
		// &models.Country{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
	}
}
