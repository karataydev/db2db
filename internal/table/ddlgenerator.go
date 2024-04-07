package table

import (
	"database/sql"
	"fmt"
	"strings"
)

type DDLGenerator interface {
	GenerateDDL(*sql.DB, string) (string, error)
}

type SQLServerDDLGenerator struct{}

func (g *SQLServerDDLGenerator) GenerateDDL(db *sql.DB, tableName string) (string, error) {
	var ddl string

	// 1. Basic Table Structure
	ddl += fmt.Sprintf("CREATE TABLE [%s] (\n", tableName)
	// 2. Get Column Details
	columnDetails, err := g.columnDetails(db, tableName)
	if err != nil {
		return "", nil
	}
	ddl += columnDetails
	// 3. Primary Key
	pk, err := g.primaryKey(db, tableName)
	if err != nil {
		return "", nil
	}
	ddl += pk
	ddl += "\n);\n\n"
	// 3. Foreign Keys
	fks, err := g.foreignKeys(db, tableName)
	if err != nil {
		return "", nil
	}
	ddl += fks
	return ddl, nil
}

func (g *SQLServerDDLGenerator) foreignKeys(db *sql.DB, tableName string) (string, error) {
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

	rows, err := db.Query(query, tableName)
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

func (g *SQLServerDDLGenerator) primaryKey(db *sql.DB, tableName string) (string, error) {
	pkStr := ""
	pkQuery := `
        SELECT c.COLUMN_NAME
        FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
        JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE c
        ON tc.CONSTRAINT_NAME = c.CONSTRAINT_NAME 
        WHERE tc.CONSTRAINT_TYPE = 'PRIMARY KEY' AND c.TABLE_NAME = ?`
	pkRows, err := db.Query(pkQuery, tableName)
	if err != nil {
		return pkStr, err
	}
	defer pkRows.Close()

	var pkCols []string
	for pkRows.Next() {
		var pkCol string
		if err := pkRows.Scan(&pkCol); err != nil {
			return "", err
		}
		pkCols = append(pkCols, pkCol)
	}

	if len(pkCols) > 0 {
		pkStr = fmt.Sprintf(",\n\tCONSTRAINT PK_%s PRIMARY KEY (%s)", tableName, strings.Join(pkCols, ", "))
	}

	return pkStr, nil
}

func (g *SQLServerDDLGenerator) columnDetails(db *sql.DB, tableName string) (string, error) {
	rows, err := db.Query(`
        SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, CHARACTER_MAXIMUM_LENGTH, 
               NUMERIC_PRECISION, NUMERIC_SCALE, COLUMN_DEFAULT
        FROM INFORMATION_SCHEMA.COLUMNS
        WHERE TABLE_NAME = ?`, tableName)

	if err != nil {
		return "", err
	}
	defer rows.Close()
	var columns []string
	for rows.Next() {
		var colName, dataType, isNullable string
		var columnDefault *string

		var maxLength, precision, scale *int // Can be null for certain data types
		if err := rows.Scan(&colName, &dataType, &isNullable, &maxLength, &precision, &scale, &columnDefault); err != nil {
			return "", err
		}
		columns = append(columns, fmt.Sprintf("\t%s %s %s", colName, dataType, g.columnTypeOptions(dataType, maxLength, precision, scale, isNullable, columnDefault)))
	}

	return strings.Join(columns, ",\n"), nil
}

func (g *SQLServerDDLGenerator) formatDefaultValue(defaultValue string) string {
	// Handle special cases (e.g., dates, system functions like GETDATE())
	if strings.Contains(defaultValue, "(") {
		return "DEFAULT " + defaultValue // Keep parentheses for function calls
	} else {
		return "DEFAULT " + fmt.Sprintf("'%s'", defaultValue) // Add quotes
	}
}

func (g *SQLServerDDLGenerator) columnTypeOptions(dataType string, maxLength, precision, scale *int, isNullable string, defaultValue *string) string {
	options := ""
	if strings.EqualFold(dataType, "BIGINT") || strings.EqualFold(dataType, "SMALLINT") || strings.EqualFold(dataType, "INT") {
		options = "" // No precision or scale needed
	} else if precision != nil && scale != nil {
		options = fmt.Sprintf("(%d, %d)", *precision, *scale)
	} else if maxLength != nil {
		options = fmt.Sprintf("(%d)", *maxLength)
	}

	// Include nullability
	if isNullable == "YES" {
		options += " NULL"
	} else {
		options += " NOT NULL"
	}

	if defaultValue != nil && *defaultValue != "" {
		options += " " + g.formatDefaultValue(*defaultValue)
	}

	return options
}
