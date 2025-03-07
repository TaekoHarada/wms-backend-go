package routes

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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})	// HTTPレスポンスとしてデータを JSON形式で返す
		return
	}
	defer rows.Close()	// defer で関数が終了する際に実行される処理を登録。データベースから取得した結果セットrowsが保持しているリソースを解放するために、rows.Close() を登録

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

// 商品取得 by ID
func GetProductByID(c *gin.Context) {
	var product struct {
		ID       int    `json:"id"`
		SKU      string `json:"sku"`
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		Location string `json:"location"`
		Category string `json:"category"`
	}

	id := c.Param("id")
	err := config.DB.QueryRow(`
		SELECT p.id, p.sku, p.name, p.quantity, p.location, c.name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = ?
	`, id).Scan(&product.ID, &product.SKU, &product.Name, &product.Quantity, &product.Location, &product.Category)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// 商品追加
func CreateProduct(c *gin.Context) {
	var product struct {
		SKU        string `json:"sku"`
		Name       string `json:"name"`
		Quantity   int    `json:"quantity"`
		Location   string `json:"location"`
		CategoryID int    `json:"category_id"`
	}

	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	result, err := config.DB.Exec(`
		INSERT INTO products (sku, name, quantity, location, category_id)
		VALUES (?, ?, ?, ?, ?)
	`, product.SKU, product.Name, product.Quantity, product.Location, product.CategoryID)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	productID, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"message": "Product added", "product_id": productID})
}

// 商品更新
func UpdateProduct(c *gin.Context) {
	var product struct {
		ID         int    `json:"id"`
		SKU        string `json:"sku"`
		Name       string `json:"name"`
		Quantity   int    `json:"quantity"`
		Location   string `json:"location"`
		CategoryID int    `json:"category_id"`
	}
	id := c.Param("id")

	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	_, err := config.DB.Exec(`
		UPDATE products
		SET sku = ?, name = ?, quantity = ?, location = ?, category_id = ?
		WHERE id = ?
	`, product.SKU, product.Name, product.Quantity, product.Location, product.CategoryID, id)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

// 商品削除
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	_, err := config.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}