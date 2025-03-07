package routes

import (
	"log"
	"net/http"

	"wms-backend-go/config"
	"wms-backend-go/models"

	"github.com/gin-gonic/gin"
)

// 在庫更新（入庫・出庫）
func UpdateStock(c *gin.Context) {
	var request struct {
		ProductID int    `json:"product_id"`
		Type      string `json:"type"` // "IN" or "OUT"
		Quantity  int    `json:"quantity"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// 現在の在庫取得
	var currentQuantity int
	err := config.DB.QueryRow("SELECT quantity FROM products WHERE id = ?", request.ProductID).Scan(&currentQuantity)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	// 在庫計算
	newQuantity := currentQuantity
	if request.Type == "IN" {
		newQuantity += request.Quantity
	} else if request.Type == "OUT" {
		if currentQuantity < request.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Insufficient stock"})
			return
		}
		newQuantity -= request.Quantity
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid transaction type"})
		return
	}

	// 在庫更新
	_, err = config.DB.Exec("UPDATE products SET quantity = ? WHERE id = ?", newQuantity, request.ProductID)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	// 在庫履歴追加
	_, err = config.DB.Exec("INSERT INTO stock_transactions (product_id, type, quantity) VALUES (?, ?, ?)",
		request.ProductID, request.Type, request.Quantity)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated", "new_quantity": newQuantity})
}

// 在庫履歴取得
func GetStockHistory(c *gin.Context) {
	rows, err := config.DB.Query(`
        SELECT st.id, st.product_id, p.name AS product_name, st.type, st.quantity, st.transaction_date
        FROM stock_transactions st
        JOIN products p ON st.product_id = p.id
        ORDER BY st.transaction_date DESC
    `)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer rows.Close()

	var history []models.StockTransaction
	for rows.Next() {
		var record models.StockTransaction
		if err := rows.Scan(&record.ID, &record.ProductID, &record.ProductName, &record.Type, &record.Quantity, &record.TransactionDate); err != nil {
			log.Println("Scan Error:", err)
			continue
		}
		history = append(history, record)
	}

	c.JSON(http.StatusOK, history)
}
