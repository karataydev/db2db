package table

import (
	"fmt"
	"log"

	"github.com/karataymarufemre/db2db/internal/connection"
)

type Creator interface {
	Create(tables []string) error
}

type SQLCreator struct {
	db           *connection.DBConnection
	ddlGenerator DDLGenerator
}

func NewSQLCreator(db *connection.DBConnection, ddlGenerator *DDLGenerator) *SQLCreator {
	return &SQLCreator{db: db, ddlGenerator: *ddlGenerator}
}

func (c *SQLCreator) Create(tables []string) error {
	q := ""
	for _, t := range tables {
		ddl, err := c.ddlGenerator.GenerateDDL(c.db.Source, t)
		if err != nil {
			log.Printf("Error generating DDL for t %s: %v\n", t, err)
			continue // Skip to the next t
		}
		q += ddl
	}
	if q != "" {
		_, err := c.db.Target.Query(q)
		if err != nil {
			return fmt.Errorf("Error executing query: %v\n", err)
		}
	}
	return nil
}
