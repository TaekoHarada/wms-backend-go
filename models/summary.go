package models

type Summary struct {
	TotalProducts  int `json:"totalProducts"`	// 総商品数
	LowStockCount  int `json:"lowStockCount"`	// 在庫切れ商品数
	RecentTransactions int `json:"recentTransactions"`
}


type StockTrend struct {
	Month    string `json:"month"`
	StockIn  int    `json:"stock_in"`
	StockOut int    `json:"stock_out"`
}
