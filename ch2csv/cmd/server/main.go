package main

import (
	"log"
	"net/http"

	"github.com/anshiq/ch2csv/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./static")

	// API routes
	api := router.Group("/api")
	{
		api.POST("/connect", handlers.ConnectHandler)
		api.GET("/tables", handlers.ListTablesHandler)
		api.GET("/columns", handlers.ListColumnsHandler)
		api.POST("/preview", handlers.PreviewDataHandler)
		api.POST("/ingest", handlers.IngestDataHandler)
	}

	// Frontend route
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
