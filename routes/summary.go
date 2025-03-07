package routes

import (
	"log"
	"net/http"

	"wms-backend-go/config"
	"wms-backend-go/models"

	"github.com/gin-gonic/gin"
)

// 在庫サマリー取得
func GetSummary(c *gin.Context) {
	var summary models.Summary
	err := config.DB.QueryRow("SELECT SUM(quantity) AS total_products FROM products").Scan(&summary.TotalProducts)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	err = config.DB.QueryRow("SELECT COUNT(*) AS low_stock_count FROM products WHERE quantity <= 5").Scan(&summary.LowStockCount)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
