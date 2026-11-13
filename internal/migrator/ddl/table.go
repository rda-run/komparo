package ddl

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/rda-run/komparo/pkg/schema"
)

func GenerateCreateTableDDL(table schema.Table) string {
	// The snapshot does not record column positions; sort by name so the
	// output is at least stable between runs.
	var cols []string
	for _, name := range slices.Sorted(maps.Keys(table.Columns)) {
		cols = append(cols, GenerateColumnDef(table.Columns[name]))
	}

	return fmt.Sprintf("CREATE TABLE %s (\n    %s\n);",
		q(table.Name),
		strings.Join(cols, ",\n    "))
}

func GenerateDropTableDDL(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", q(tableName))
}

func GenerateAddColumnDDL(tableName string, col schema.Column) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s;", q(tableName), GenerateColumnDef(col))
}

func GenerateDropColumnDDL(tableName, colName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s;", q(tableName), q(colName))
}

func GenerateAlterColumnTypeDDL(tableName, colName string, col schema.Column) string {
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;",
		q(tableName), q(colName), formatColumnType(col))
}

func GenerateAlterColumnNullableDDL(tableName, colName string, nullable bool) string {
	if nullable {
		return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s DROP NOT NULL;", q(tableName), q(colName))
	}
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;", q(tableName), q(colName))
}

func GenerateAlterColumnDefaultDDL(tableName, colName string, defaultValue *string) string {
	if defaultValue == nil {
		return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s DROP DEFAULT;", q(tableName), q(colName))
	}
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET DEFAULT %s;", q(tableName), q(colName), *defaultValue)
}

func GenerateColumnDef(col schema.Column) string {
	parts := []string{q(col.Name), formatColumnType(col)}

	if col.Collation != nil {
		parts = append(parts, "COLLATE", q(*col.Collation))
	}

	if !col.Nullable {
		parts = append(parts, "NOT NULL")
	}

	if col.Default != nil {
		parts = append(parts, "DEFAULT", *col.Default)
	}

	return strings.Join(parts, " ")
}

func formatColumnType(col schema.Column) string {
	typeUpper := strings.ToUpper(col.Type)

	if col.MaxLen != nil && *col.MaxLen > 0 {
		if typeAcceptsLength(typeUpper) {
			return fmt.Sprintf("%s(%d)", col.Type, *col.MaxLen)
		}
	}

	if col.NumericPrec != nil && col.NumericScale != nil {
		if typeAcceptsPrecisionScale(typeUpper) {
			return fmt.Sprintf("%s(%d,%d)", col.Type, *col.NumericPrec, *col.NumericScale)
		}
	}

	if col.NumericPrec != nil {
		if typeAcceptsPrecisionScale(typeUpper) {
			return fmt.Sprintf("%s(%d)", col.Type, *col.NumericPrec)
		}
	}

	return col.Type
}

func typeAcceptsLength(typeUpper string) bool {
	switch typeUpper {
	case "VARCHAR", "CHAR", "CHARACTER":
		return true
	}
	return false
}

func typeAcceptsPrecisionScale(typeUpper string) bool {
	switch typeUpper {
	case "NUMERIC", "DECIMAL":
		return true
	}
	return false
}
