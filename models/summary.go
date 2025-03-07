package models

type Summary struct {
	TotalProducts  int `json:"total_products"`	// 総商品数
	LowStockCount  int `json:"low_stock_count"`	// 在庫切れ商品数
}
