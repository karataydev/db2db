package table

import (
	"database/sql"
	"fmt"
	"strings"
)

type DDLGenerator interface {
	GenerateDDL(string) (string, error)
}

type SQLServerDDLGenerator struct {
	db *sql.DB
	c  ColumnDetailService
}

func NewSQLServerDDLGenerator(db *sql.DB, c ColumnDetailService) *SQLServerDDLGenerator {
	return &SQLServerDDLGenerator{db: db, c: c}
}

func (g *SQLServerDDLGenerator) GenerateDDL(tableName string) (string, error) {
	var ddl string

	// 1. Basic Table Structure
	ddl += fmt.Sprintf("CREATE TABLE [%s] (\n", tableName)
	// 2. Get Column Details
	cols, err := g.c.ColumnDetail(tableName)
	if err != nil {
		return "", err
	}
	columnDetails, err := g.columnDetails(cols)
	if err != nil {
		return "", nil
	}
	ddl += columnDetails
	// 3. Primary Key
	pk := g.primaryKey(tableName, cols)
	ddl += pk
	ddl += "\n);\n\n"
	// 3. Foreign Keys
	fks, err := g.foreignKeys(tableName)
	if err != nil {
		return "", nil
	}
	ddl += fks
	fmt.Println(ddl)
	return ddl, nil
}

func (g *SQLServerDDLGenerator) foreignKeys(tableName string) (string, error) {
	foreignKeys := ""
	query := `
	SELECT 
		rc.CONSTRAINT_NAME,
		FKCU.TABLE_NAME, 
		FKCU.COLUMN_NAME AS Columns, 
		RCU.TABLE_NAME AS ReferencedTable,
		RCU.COLUMN_NAME AS ReferencedColumns
	FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc
	JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE FKCU ON rc.CONSTRAINT_NAME = FKCU.CONSTRAINT_NAME
	JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE RCU ON rc.UNIQUE_CONSTRAINT_NAME = RCU.CONSTRAINT_NAME
	WHERE FKCU.TABLE_NAME =  ?
    `

	rows, err := g.db.Query(query, tableName)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var ConstraintName, TableName, Columns, ReferencedTable, ReferencedColumns string
	for rows.Next() {
		if err := rows.Scan(&ConstraintName, &TableName, &Columns, &ReferencedTable, &ReferencedColumns); err != nil {
			return "", err
		}
		foreignKeys = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s);\n",
			TableName, ConstraintName, Columns, ReferencedTable, ReferencedColumns)
	}

	return foreignKeys, nil
}

func (g *SQLServerDDLGenerator) primaryKey(tableName string, cols []Column) string {
	var pkCols []string
	for _, col := range cols {
		if col.IsPrimaryKey {
			pkCols = append(pkCols, col.ColName)
		}
	}

	if len(pkCols) > 0 {
		return fmt.Sprintf(",\n\tCONSTRAINT PK_%s PRIMARY KEY (%s)", tableName, strings.Join(pkCols, ", "))
	}

	return ""
}

func (g *SQLServerDDLGenerator) columnDetails(columns []Column) (string, error) {
	var columnsStr []string
	for _, col := range columns {
		columnsStr = append(columnsStr, fmt.Sprintf("\t%s %s %s", col.ColName, col.DataType, g.columnTypeOptions(col)))
	}
	return strings.Join(columnsStr, ",\n"), nil
}

func (g *SQLServerDDLGenerator) formatDefaultValue(defaultValue string) string {
	// Handle special cases (e.g., dates, system functions like GETDATE())
	if strings.Contains(defaultValue, "(") {
		return "DEFAULT " + defaultValue // Keep parentheses for function calls
	} else {
		return "DEFAULT " + fmt.Sprintf("'%s'", defaultValue) // Add quotes
	}
}

func (g *SQLServerDDLGenerator) columnTypeOptions(col Column) string {
	options := ""
	if col.isIdentity {
		options = "IDENTITY(1,1)" // Include IDENTITY
	} else if strings.EqualFold(col.DataType, "BIGINT") || strings.EqualFold(col.DataType, "SMALLINT") || strings.EqualFold(col.DataType, "INT") {
		options = "" // No precision or scale needed
	} else if col.precision != nil && col.scale != nil {
		options = fmt.Sprintf("(%d, %d)", *col.precision, *col.scale)
	} else if col.maxLength != nil {
		options = fmt.Sprintf("(%d)", *col.maxLength)
	}

	// Include nullability
	if col.isNullable == "YES" {
		options += " NULL"
	} else {
		options += " NOT NULL"
	}

	if col.columnDefault != nil && *col.columnDefault != "" {
		options += " " + g.formatDefaultValue(*col.columnDefault)
	}

	return options
}
