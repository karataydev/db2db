package table

import (
	"database/sql"
)

type Column struct {
	ColName       string
	DataType      string
	isNullable    string
	columnDefault *string
	maxLength     *int
	precision     *int
	scale         *int
	IsPrimaryKey  bool
	isIdentity    bool
}

type ColumnDetailService interface {
	ColumnDetail(tableName string) ([]Column, error)
}

type SQLServerColumnDetailService struct {
	db *sql.DB
}

func NewSQLServerColumnDetailService(db *sql.DB) *SQLServerColumnDetailService {
	return &SQLServerColumnDetailService{db: db}
}

func (c *SQLServerColumnDetailService) ColumnDetail(tableName string) ([]Column, error) {
	var columns []Column

	rows, err := c.db.Query(`
			SELECT c.COLUMN_NAME, c.DATA_TYPE, c.IS_NULLABLE, c.CHARACTER_MAXIMUM_LENGTH, 
				   c.NUMERIC_PRECISION, c.NUMERIC_SCALE, c.COLUMN_DEFAULT,
				   IIF(tc.CONSTRAINT_TYPE = 'PRIMARY KEY', 1, 0) AS IS_PRIMARY,
				   COLUMNPROPERTY(OBJECT_ID(c.TABLE_SCHEMA+'.'+c.TABLE_NAME), c.COLUMN_NAME, 'IsIdentity')
			FROM INFORMATION_SCHEMA.COLUMNS c
			LEFT OUTER JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kc ON kc.COLUMN_NAME = c.COLUMN_NAME AND kc.TABLE_NAME = c.TABLE_NAME
			LEFT OUTER JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc ON tc.CONSTRAINT_NAME = kc.CONSTRAINT_NAME
			WHERE c.TABLE_NAME = ?`, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		col := &Column{}
		if err := rows.Scan(&col.ColName, &col.DataType, &col.isNullable, &col.maxLength, &col.precision, &col.scale, &col.columnDefault, &col.IsPrimaryKey, &col.isIdentity); err != nil {
			return nil, err
		}
		columns = append(columns, *col)
	}
	return columns, nil
}
