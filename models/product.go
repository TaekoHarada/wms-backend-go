package models

type Product struct {
	ID         int    `json:"id"`
	SKU        string `json:"sku"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	Location   string `json:"location"`
	Category   string `json:"category"`
	CategoryID int    `json:"category_id"`
}
