package services

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/anshiq/ch2csv/internal/models"
)

func TestClickHouseConnection(config models.ClickHouseConfig) error {
	conn, err := createClickHouseConnection(config)
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.Ping(context.Background())
}

func ListClickHouseTables(config models.ClickHouseConfig) ([]string, error) {
	conn, err := createClickHouseConnection(config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := "SELECT name FROM system.tables WHERE database = ?"
	rows, err := conn.Query(context.Background(), query, config.Database)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func ListClickHouseColumns(config models.ClickHouseConfig, table string) ([]models.ColumnInfo, error) {
	conn, err := createClickHouseConnection(config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := fmt.Sprintf("SELECT name, type FROM system.columns WHERE database = ? AND table = ?")
	rows, err := conn.Query(context.Background(), query, config.Database, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []models.ColumnInfo
	for rows.Next() {
		var col models.ColumnInfo
		if err := rows.Scan(&col.Name, &col.Type); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, nil
}

func PreviewClickHouseData(config models.ClickHouseConfig, table string, columns []string, limit int) ([]map[string]interface{}, error) {
	conn, err := createClickHouseConnection(config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := fmt.Sprintf("SELECT %s FROM %s LIMIT %d", formatColumns(columns), table, limit)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMap(rows)
}

func ClickHouseToCSV(config models.ClickHouseConfig, table string, columns []string, filePath string, delimiter rune) (int, error) {
	conn, err := createClickHouseConnection(config)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	query := fmt.Sprintf("SELECT %s FROM %s", formatColumns(columns), table)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	return writeToCSV(rows, filePath, delimiter)
}

func createClickHouseConnection(config models.ClickHouseConfig) (driver.Conn, error) {
	fmt.Print(config)
	if config.JWTToken == "" {

		return nil, fmt.Errorf("No password as jwt provided")
	}

	opts := &clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", config.Host, config.Port)},
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.User,
			Password: config.JWTToken,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
	}

	if config.Secure {
		opts.TLS = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return clickhouse.Open(opts)
}

// Helper functions
func formatColumns(columns []string) string {
	if len(columns) == 0 {
		return "*"
	}

	var formatted string
	for i, col := range columns {
		if i > 0 {
			formatted += ", "
		}
		formatted += fmt.Sprintf("`%s`", col)
	}
	return formatted
}

func scanRowsToMap(rows driver.Rows) ([]map[string]interface{}, error) {
	columns := rows.Columns()
	var result []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}
		result = append(result, rowMap)
	}

	return result, nil
}
