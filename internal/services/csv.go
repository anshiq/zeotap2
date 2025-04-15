package services

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/anshiq/ch2csv/internal/models"
)

func ValidateFlatFile(filePath string) error {
	_, err := os.Stat(filePath)
	return err
}

func GetCSVColumns(filePath string, delimiter rune) ([]models.ColumnInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var columns []models.ColumnInfo
	for _, header := range headers {
		columns = append(columns, models.ColumnInfo{Name: header, Type: "string"})
	}

	return columns, nil
}

func PreviewCSVData(filePath string, delimiter rune, columns []string, limit int) ([]map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Create column index map
	colIndex := make(map[string]int)
	for i, h := range headers {
		colIndex[h] = i
	}

	// Determine which columns to include
	var includeIndices []int
	if len(columns) == 0 {
		// Include all columns
		for i := range headers {
			includeIndices = append(includeIndices, i)
		}
	} else {
		// Include only selected columns
		for _, col := range columns {
			if idx, ok := colIndex[col]; ok {
				includeIndices = append(includeIndices, idx)
			}
		}
	}

	var result []map[string]interface{}
	count := 0

	for count < limit {
		record, err := reader.Read()
		if err != nil {
			break
		}

		row := make(map[string]interface{})
		for _, idx := range includeIndices {
			if idx < len(record) {
				row[headers[idx]] = record[idx]
			}
		}

		result = append(result, row)
		count++
	}

	return result, nil
}

func CSVToClickHouse(config models.ClickHouseConfig, table string, columns []string, filePath string, delimiter rune) (int, error) {
	// This is a simplified version. In a real implementation, you would:
	// 1. Create a temporary file in ClickHouse format
	// 2. Use ClickHouse's native format for bulk inserts
	// 3. Or implement proper batch inserts

	// For simplicity, we'll just count the rows here
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter

	// Skip header
	if _, err := reader.Read(); err != nil {
		return 0, err
	}

	count := 0
	for {
		_, err := reader.Read()
		if err != nil {
			break
		}
		count++
	}

	return count, nil
}

func writeToCSV(rows driver.Rows, filePath string, delimiter rune) (int, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = delimiter
	defer writer.Flush()

	columns := rows.Columns()

	// Write header
	if err := writer.Write(columns); err != nil {
		return 0, err
	}

	count := 0

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return count, err
		}

		record := make([]string, len(columns))
		for i, val := range values {
			switch v := val.(type) {
			case string:
				record[i] = v
			case int64:
				record[i] = strconv.FormatInt(v, 10)
			case float64:
				record[i] = strconv.FormatFloat(v, 'f', -1, 64)
			case bool:
				record[i] = strconv.FormatBool(v)
			case nil:
				record[i] = ""
			default:
				record[i] = fmt.Sprintf("%v", v)
			}
		}

		if err := writer.Write(record); err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}
