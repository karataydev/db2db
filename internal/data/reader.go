package data

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/karataymarufemre/db2db/internal/table"
	"github.com/shopspring/decimal"
)

type Reader interface {
	Read(tableName string, page int) (string, error)
}

type SQLServerReader struct {
	db    *sql.DB
	c     table.ColumnDetailService
	limit int
}

func NewSQLServerReader(db *sql.DB, c table.ColumnDetailService) *SQLServerReader {
	return &SQLServerReader{db: db, c: c, limit: 2}
}

func (s *SQLServerReader) Read(tableName string, page int) (string, error) {
	cols, err := s.c.ColumnDetail(tableName)
	if err != nil {
		return "", err
	}
	var dmlScript strings.Builder
	dmlScript.WriteString(fmt.Sprintf("INSERT INTO %s (", tableName))

	var columnNames []string
	var primaryKey string
	for _, v := range cols {
		columnNames = append(columnNames, v.ColName)
		if v.IsPrimaryKey {
			primaryKey = primaryKey + v.ColName + ","
		}
	}
	primaryKey = strings.TrimSuffix(primaryKey, ",")
	dmlScript.WriteString(strings.Join(columnNames, ", "))
	dmlScript.WriteString(") VALUES \n")
	offset := page * s.limit
	rows, err := s.db.Query(fmt.Sprintf("SELECT * FROM %s ORDER BY %s OFFSET %d ROWS FETCH NEXT %d ROWS ONLY;", tableName, primaryKey, offset, s.limit))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var values []interface{}
	valueExists := false

	for rows.Next() {
		valueExists = true
		values = values[:0] // Reset values for each row

		// Create a slice of interface{} to scan, matching the number of columns
		values = make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// Read the row into the interface slice
		err = rows.Scan(scanArgs...)
		if err != nil {
			return "", err
		}

		dmlScript.WriteString("(")
		for i := range cols {
			dmlScript.WriteString(formatValue(values[i]))
			if i < len(cols)-1 {
				dmlScript.WriteString(", ")
			}
		}
		dmlScript.WriteString("),\n")
	}

	if !valueExists {
		return "", nil
	}

	return strings.TrimSuffix(dmlScript.String(), ",\n") + ";\n", nil
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'" // Escape single quotes
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case []uint8: // Numeric/Decimal represented as []uint8
		strValue := string(v) // Convert to string first
		if decVal, err := decimal.NewFromString(strValue); err == nil {
			return decVal.String()
		} else {
			// Handle the error if it's not a valid decimal either
			return "NULL" // Or handle the error differently
		}
	case bool:
		if v {
			return "1"
		} else {
			return "0"
		}
	case time.Time:
		// Important: Adapt the format string based on your database's expected format
		return "'" + v.Format("2006-01-02 15:04:05") + "'"
	case nil: // Handle NULL values
		return "NULL"
	default:
		return "NULL" // Placeholder for potentially unsupported types
	}
}
