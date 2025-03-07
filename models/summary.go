package models

type Summary struct {
	TotalProducts  int `json:"total_products"`
	LowStockCount  int `json:"low_stock_count"`
}
