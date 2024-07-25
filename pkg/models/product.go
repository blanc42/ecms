package models

import (
	"gorm.io/gorm"
)

// Admin model
type Admin struct {
	gorm.Model
	Username string `gorm:"type:varchar(255);not null;uniqueIndex" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Email    string `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Stores   []Store
}

// Store model
type Store struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null;index" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	AdminID     uint   `gorm:"not null" json:"admin_id"`
	Admin       *Admin `gorm:"foreignKey:AdminID"`
	Categories  []Category
	Products    []Product  `json:"-"`
	Customers   []Customer `json:"-"`
	Orders      []Order    `json:"-"`
}

// Category model
type Category struct {
	gorm.Model
	Name             string    `gorm:"type:varchar(255);not null;index" json:"name"`
	Description      string    `gorm:"type:text" json:"description"`
	StoreID          uint      `gorm:"not null;index" json:"store_id"`
	Store            *Store    `gorm:"foreignKey:StoreID"`
	ParentCategoryID *uint     `json:"parent_category_id,omitempty"`
	ParentCategory   *Category `gorm:"foreignKey:ParentCategoryID"`
	// Subcategories    []*Category `gorm:"foreignKey:ParentCategoryID"`
	Subcategories []*Category `gorm:"-"`
	Products      []Product   `json:"-"`
	Variants      []Variant   `json:"-"`
}

// Product model
type Product struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(255);not null;index" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Rating      float32   `json:"rating"`
	IsFeatured  bool      `json:"is_featured"`
	IsArchived  bool      `json:"is_archived"`
	HasVariants bool      `json:"has_variants"`
	CategoryID  uint      `gorm:"not null;index" json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID"`
	StoreID     uint      `gorm:"not null;index" json:"store_id"`
	Store       *Store    `gorm:"foreignKey:StoreID"`
	Items       []ProductItem
	Images      []ProductImage
}

// Variant model
type Variant struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(255);not null;index" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Weight      int       `gorm:"not null;gt:0" json:"weight"`
	CategoryID  uint      `gorm:"not null;index" json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID"`
	Options     []VariantOption
}

// VariantOption model
type VariantOption struct {
	gorm.Model
	Value       string   `gorm:"type:varchar(255);not null" json:"value"`
	Description string   `gorm:"type:text" json:"description"`
	Weight      int      `json:"weight"`
	VariantID   uint     `gorm:"not null;index" json:"variant_id"`
	Variant     *Variant `gorm:"foreignKey:VariantID"`
}

// ProductItem model
type ProductItem struct {
	gorm.Model
	ProductID       uint     `gorm:"not null;index" json:"product_id"`
	Product         *Product `gorm:"foreignKey:ProductID"`
	SKU             string   `gorm:"type:varchar(255);not null;uniqueIndex" json:"sku"`
	Quantity        int      `gorm:"not null" json:"quantity"`
	Price           float64  `gorm:"type:float;not null" json:"price"`
	DiscountedPrice float64  `gorm:"type:float" json:"discounted_price,omitempty"`
}

// ProductImage model
type ProductImage struct {
	gorm.Model
	ProductID uint     `gorm:"not null;index" json:"product_id"`
	Product   *Product `gorm:"foreignKey:ProductID"`
	ImageURL  string   `gorm:"type:varchar(255); not null" json:"image_url"`
}

// Customer model
type Customer struct {
	gorm.Model
	Username  string   `gorm:"type:varchar(255);not null;uniqueIndex" json:"username"`
	Password  string   `gorm:"type:varchar(255); not null" json:"-"`
	Email     string   `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	StoreID   uint     `gorm:"not null;index" json:"store_id"`
	Store     *Store   `gorm:"foreignKey:StoreID"`
	AddressID uint     `json:"address_id,omitempty"`
	Address   *Address `gorm:"foreignKey:AddressID"`
	Carts     []Cart
	Orders    []Order
}

// Order model
type Order struct {
	gorm.Model
	OrderNumber   string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"order_number"`
	PaymentStatus string    `gorm:"type:varchar(50); not null" json:"payment_status"`
	OrderStatus   string    `gorm:"type:varchar(50); not null" json:"order_status"`
	StoreID       uint      `gorm:"not null;index" json:"store_id"`
	Store         *Store    `gorm:"foreignKey:StoreID"`
	CustomerID    uint      `gorm:"not null;index" json:"customer_id"`
	Customer      *Customer `gorm:"foreignKey:CustomerID"`
	OrderItems    []OrderItem
}

// OrderItem model
type OrderItem struct {
	gorm.Model
	ProductItemID uint         `gorm:"not null;index" json:"product_item_id"`
	ProductItem   *ProductItem `gorm:"foreignKey:ProductItemID"`
	Quantity      int          `gorm:"not null" json:"quantity"`
	OrderID       uint         `gorm:"not null;index" json:"order_id"`
	Order         *Order       `gorm:"foreignKey:OrderID"`
}

// Cart model
type Cart struct {
	gorm.Model
	CustomerID uint      `gorm:"not null;index" json:"customer_id"`
	Customer   *Customer `gorm:"foreignKey:CustomerID"`
	CartItems  []CartItem
}

// CartItem model
type CartItem struct {
	gorm.Model
	ProductItemID uint         `gorm:"not null;index" json:"product_item_id"`
	ProductItem   *ProductItem `gorm:"foreignKey:ProductItemID"`
	Quantity      int          `gorm:"not null" json:"quantity"`
	CartID        uint         `gorm:"not null;index" json:"cart_id"`
	Cart          *Cart        `gorm:"foreignKey:CartID"`
}

// Address model
type Address struct {
	gorm.Model
	City      string   `gorm:"type:varchar(255);not null" json:"city"`
	Pincode   string   `gorm:"type:varchar(255);not null" json:"pincode"`
	CountryID uint     `gorm:"not null;index" json:"country_id"`
	Country   *Country `gorm:"foreignKey:CountryID"`
}

// Country model
type Country struct {
	gorm.Model
	Country string `gorm:"type:varchar(255);not null" json:"country"`
}
