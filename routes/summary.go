package routes

import (
	"database/sql"
	"log"
	"net/http"

	"wms-backend-go/config"
	"wms-backend-go/models"

	"github.com/gin-gonic/gin"
)

// 在庫サマリー取得 API
func GetSummary(c *gin.Context) {
	var summary models.Summary

	err := config.DB.QueryRow("SELECT IFNULL(SUM(quantity), 0) AS totalProducts FROM products").Scan(&summary.TotalProducts)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	err = config.DB.QueryRow("SELECT COUNT(*) AS lowStockCount FROM products WHERE quantity <= 5").Scan(&summary.LowStockCount)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	err = config.DB.QueryRow("SELECT COUNT(*) AS recentTransactions FROM stock_transactions WHERE transaction_date >= NOW() - INTERVAL 5 DAY").Scan(&summary.RecentTransactions)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// 在庫切れ商品の取得 API
func GetLowStock(c *gin.Context) {
	rows, err := config.DB.Query(`
								SELECT p.id, p.sku, p.name, p.quantity, p.location, c.name, p.category_id 
								FROM products p 
								JOIN categories c ON p.category_id = c.id
								WHERE p.quantity <= 5 
								ORDER BY p.quantity ASC
								`)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Quantity,
			&product.Location,
			&product.Category,
			&product.CategoryID,
		); err != nil {
			log.Println("Database Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

// 最近の入庫・出庫履歴取得 API
func GetRecentStock(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT st.id, p.name AS product_name, st.type, st.quantity, st.transaction_date
		FROM stock_transactions st
		JOIN products p ON st.product_id = p.id
		ORDER BY st.transaction_date DESC
		LIMIT 5
	`)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer rows.Close()

	var transactions []models.StockTransaction
	for rows.Next() {
		var transaction models.StockTransaction
		if err := rows.Scan(&transaction.ID, &transaction.ProductName, &transaction.Type, &transaction.Quantity, &transaction.TransactionDate); err != nil {
			log.Println("Database Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
			return
		}
		transactions = append(transactions, transaction)
	}

	c.JSON(http.StatusOK, transactions)
}

// 月ごとの在庫推移取得 API
func GetStockTrends(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT 
			DATE_FORMAT(transaction_date, '%Y-%m') AS month,
			IFNULL(SUM(CASE WHEN type = 'IN' THEN quantity ELSE 0 END), 0) AS stock_in,
			IFNULL(SUM(CASE WHEN type = 'OUT' THEN quantity ELSE 0 END), 0) AS stock_out
		FROM stock_transactions
		WHERE transaction_date >= DATE_SUB(NOW(), INTERVAL 6 MONTH)
		GROUP BY month
		ORDER BY month ASC
	`)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer rows.Close()

	var trends []models.StockTrend
	for rows.Next() {
		var trend models.StockTrend
		if err := rows.Scan(&trend.Month, &trend.StockIn, &trend.StockOut); err != nil {
			log.Println("Database Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
			return
		}
		trends = append(trends, trend)
	}

	c.JSON(http.StatusOK, trends)
}
