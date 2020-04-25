package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func main() {
	db := setupDatabase(true)
	createMockData(db)

	// Read
	var query RoomModel
	db.First(&query, 1) // find the first product with id 1
	fmt.Println(query)

	// // Update
	// db.Model(&product).Update("Price", 2000)

	// db.Delete(&product)

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
