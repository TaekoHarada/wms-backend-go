package models

import "time"

type StockTransaction struct {
	ID              int       `json:"id"`
	ProductID       int       `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductSKU      string    `json:"sku"`
	Type            string    `json:"type"` // "IN" or "OUT"
	Quantity        int       `json:"quantity"`
	TransactionDate time.Time `json:"transaction_date"`
}
