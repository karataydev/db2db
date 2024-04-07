package table

import (
	"database/sql"
)

type TableNameReader interface {
	GetTableNames(db *sql.DB) ([]string, error)
}

type SQLServerTableNameReader struct{}

func (s *SQLServerTableNameReader) GetTableNames(db *sql.DB) (tableNames []string, err error) {
	// SQL Server query to get table names
	rows, err := db.Query(`
		SELECT TABLE_NAME
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_TYPE = 'BASE TABLE' -- Exclude views
	`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		dErr := rows.Close()
		if dErr != nil && err == nil {
			err = dErr
		}
	}(rows)

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tableNames = append(tableNames, tableName)
	}

	return tableNames, nil
}
