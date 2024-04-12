package data

import (
	"database/sql"
	"fmt"
	"log"
)

type Writer interface {
	Write(query string, tableName string) error
}

type SQLServerWriter struct {
	db *sql.DB
}

func NewSQLServerWriter(db *sql.DB) *SQLServerWriter {
	return &SQLServerWriter{db: db}
}

func (s *SQLServerWriter) Write(query string, tableName string) error {
	dml := fmt.Sprintf("SET IDENTITY_INSERT %s ON;\n", tableName) + query + fmt.Sprintf("\nSET IDENTITY_INSERT %s OFF;\n", tableName)
	_, err := s.db.Query(dml)
	if err != nil {
		return fmt.Errorf("Error executing query: %v\n", err)
	}
	log.Println(dml)
	return nil
}
