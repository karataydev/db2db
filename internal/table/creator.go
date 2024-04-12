package table

import (
	"database/sql"
	"fmt"
	"log"
)

type Creator interface {
	CreateTables(tables []string) error
}

type SQLCreator struct {
	targetDB     *sql.DB
	ddlGenerator DDLGenerator
}

func NewSQLCreator(targetDB *sql.DB, ddlGenerator *DDLGenerator) *SQLCreator {
	return &SQLCreator{targetDB: targetDB, ddlGenerator: *ddlGenerator}
}

func (c *SQLCreator) CreateTables(tables []string) error {
	q := ""
	for _, t := range tables {
		ddl, err := c.ddlGenerator.GenerateDDL(t)
		if err != nil {
			log.Printf("Error generating DDL for t %s: %v\n", t, err)
			continue // Skip to the next t
		}
		q += ddl
	}
	if q != "" {
		_, err := c.targetDB.Query(q)
		if err != nil {
			return fmt.Errorf("Error executing query: %v\n", err)
		}
	}
	return nil
}
