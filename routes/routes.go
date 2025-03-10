package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	router.GET("/products", GetAllProducts)
	router.GET("/products/:id", GetProductByID)
	router.POST("/products", CreateProduct)
	router.PUT("/products/:id", UpdateProduct)
	router.DELETE("/products/:id", DeleteProduct)

	router.GET("/categories", GetAllCategories)
	router.POST("/categories", CreateCategory)

	router.POST("/stock/update", UpdateStock)
	router.GET("/stock/history", GetStockHistory)

	router.GET("/summary", GetSummary)
	router.GET("/summary/low-stock", GetLowStock)
	router.GET("/summary/recent-stock", GetRecentStock)
	router.GET("/summary/stock-trends", GetStockTrends)
}
