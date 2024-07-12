package dto

// CartItemResponse represents the response for a cart item
type CartItemResponse struct {
	ID        uint          `json:"id"`
	ProductID uint          `json:"productId"`
	Product   ProductDetail `json:"product"`
	Quantity  uint          `json:"quantity"`
}

// ProductDetail represents the detailed information of a product in the cart
type ProductDetail struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Photo       string  `json:"photo"`
}

type UpdateQuantity struct {
	Quantity uint `json:"quantity"`
}
