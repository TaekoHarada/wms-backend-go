package routes // 統一！

import (
	"log"
	"net/http"

	"wms-backend-go/config"

	"github.com/gin-gonic/gin"
)

// 商品一覧取得
func GetAllProducts(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT p.id, p.sku, p.name, p.quantity, p.location, c.name AS category, c.id AS category_id
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
	`)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer rows.Close()

	var products []map[string]interface{}
	for rows.Next() {
		var id, categoryID, quantity int
		var sku, name, location, category string

		err := rows.Scan(&id, &sku, &name, &quantity, &location, &category, &categoryID)
		if err != nil {
			log.Println("Scan Error:", err)
			continue
		}

		product := map[string]interface{}{
			"id":          id,
			"sku":         sku,
			"name":        name,
			"quantity":    quantity,
			"location":    location,
			"category":    category,
			"category_id": categoryID,
		}
		products = append(products, product)
	}
	c.JSON(http.StatusOK, products)
}
