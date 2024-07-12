package dto

// ProductRequest represents the request body for creating or updating a product
type ProductRequest struct {
	Name        string  `form:"name" json:"name" binding:"required"`
	Description string  `form:"description" json:"description" binding:"required"`
	Price       float64 `form:"price" json:"price" binding:"required"`
	Photo       string  `form:"photo" json:"photo"`
}

// ProductResponse represents the response body for a product
type ProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Photo       string  `json:"photo"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// AddToCartRequest represents the request body for adding a product to the cart
type AddToCartRequest struct {
	ProductID uint `json:"productId" binding:"required"`
	Quantity  uint `json:"quantity" binding:"required"`
}
