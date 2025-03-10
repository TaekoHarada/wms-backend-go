package routes

import (
	"log"
	"net/http"

	"wms-backend-go/config"

	"github.com/gin-gonic/gin"
)

// カテゴリ一覧取得
func GetAllCategories(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, name FROM categories ORDER BY name ASC")
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer rows.Close()

	var categories []map[string]interface{}
	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			log.Println("Scan Error:", err)
			continue
		}

		categories = append(categories, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}
	c.JSON(http.StatusOK, categories)
}

// カテゴリ追加
func CreateCategory(c *gin.Context) {
	var request struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	_, err := config.DB.Exec("INSERT INTO categories (name) VALUES (?)", request.Name)
	if err != nil {
		log.Println("Database Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created"})
}