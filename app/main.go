package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.GET("/models", func(c *gin.Context) {
		c.String(http.StatusOK, queryRoomModels())
	})

	r.GET("/models/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, queryRoomModel(id))
	})

	r.GET("/sensors", func(c *gin.Context) {
		c.String(http.StatusOK, querySensors())
	})

	r.GET("/sensors/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, querySensor(id))
	})

	return r
}

func main() {
	setupDatabase(true)
	createMockData()

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
