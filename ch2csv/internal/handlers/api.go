package handlers

import (
	"fmt"
	"net/http"

	"github.com/anshiq/ch2csv/internal/models"
	"github.com/anshiq/ch2csv/internal/services"
	"github.com/gin-gonic/gin"
)

func ConnectHandler(c *gin.Context) {
	var req models.ConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.SourceType == "clickhouse" {
		if req.ClickHouseConfig.Host == "" || req.ClickHouseConfig.Database == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required ClickHouse connection parameters"})
			return
		}
	} else if req.SourceType == "flatfile" {
		if req.FlatFileConfig.FilePath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file path"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source type"})
		return
	}

	// Test connection
	var err error
	if req.SourceType == "clickhouse" {
		err = services.TestClickHouseConnection(req.ClickHouseConfig)
	} else {
		err = services.ValidateFlatFile(req.FlatFileConfig.FilePath)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection successful"})
}

func ListTablesHandler(c *gin.Context) {
	sourceType := c.Query("sourceType")

	if sourceType == "clickhouse" {
		config := models.ClickHouseConfig{
			Host:     c.Query("host"),
			Port:     c.Query("port"),
			Database: c.Query("database"),
			User:     c.Query("user"),
			JWTToken: c.Query("jwtToken"),
			Secure:   c.Query("secure") == "true",
		}

		tables, err := services.ListClickHouseTables(config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tables": tables})
	} else {
		c.JSON(http.StatusOK, gin.H{"tables": []string{"file_data"}})
	}
}

func ListColumnsHandler(c *gin.Context) {
	sourceType := c.Query("sourceType")
	table := c.Query("table")

	if sourceType == "clickhouse" {
		config := models.ClickHouseConfig{
			Host:     c.Query("host"),
			Port:     c.Query("port"),
			Database: c.Query("database"),
			User:     c.Query("user"),
			JWTToken: c.Query("jwtToken"),
			Secure:   c.Query("secure") == "true",
		}

		columns, err := services.ListClickHouseColumns(config, table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"columns": columns})
	} else {
		filePath := c.Query("filePath")
		delimiter := c.Query("delimiter")
		if delimiter == "" {
			delimiter = ","
		}

		columns, err := services.GetCSVColumns(filePath, rune(delimiter[0]))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"columns": columns})
	}
}

func PreviewDataHandler(c *gin.Context) {
	var req models.PreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data []map[string]interface{}
	var err error

	if req.SourceType == "clickhouse" {
		data, err = services.PreviewClickHouseData(req.ClickHouseConfig, req.Table, req.Columns, 100)
	} else {
		data, err = services.PreviewCSVData(req.FlatFileConfig.FilePath, rune(req.FlatFileConfig.Delimiter[0]), req.Columns, 100)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func IngestDataHandler(c *gin.Context) {
	var req models.IngestRequest
	config := models.ClickHouseConfig{
		Host:     c.Query("host"),
		Port:     c.Query("port"),
		Database: c.Query("database"),
		User:     c.Query("user"),
		JWTToken: c.Query("jwtToken"),
		Secure:   c.Query("secure") == "true",
	}
	req.ClickHouseConfig = config

	fmt.Print(req.ClickHouseConfig)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var count int
	var err error

	if req.Direction == "clickhouse_to_flatfile" {
		count, err = services.ClickHouseToCSV(
			req.ClickHouseConfig,
			req.Table,
			req.Columns,
			req.FlatFileConfig.FilePath,
			rune(req.FlatFileConfig.Delimiter[0]),
		)
	} else {
		count, err = services.CSVToClickHouse(
			req.ClickHouseConfig,
			req.Table,
			req.Columns,
			req.FlatFileConfig.FilePath,
			rune(req.FlatFileConfig.Delimiter[0]),
		)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
