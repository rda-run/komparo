package ddl

import (
	"fmt"

	"github.com/rda-run/komparo/pkg/schema"
)

func GenerateAddConstraintDDL(c schema.Constraint) string {
	return fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s;",
		c.TableName, c.Name, c.Definition)
}

func GenerateDropConstraintDDL(tableName, conName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s;",
		tableName, conName)
}

func GetConstraintTypeDescription(conType string) string {
	switch conType {
	case "p":
		return "PRIMARY KEY"
	case "f":
		return "FOREIGN KEY"
	case "u":
		return "UNIQUE"
	case "c":
		return "CHECK"
	case "x":
		return "EXCLUDE"
	default:
		return conType
	}
}
