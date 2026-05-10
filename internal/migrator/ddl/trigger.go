package ddl

import (
	"fmt"

	"github.com/rda-run/komparo/pkg/schema"
)

func GenerateCreateTriggerDDL(t schema.Trigger) string {
	return t.Definition + ";"
}

func GenerateDropTriggerDDL(tableName, triggerName string) string {
	return fmt.Sprintf("DROP TRIGGER IF EXISTS %s ON %s;",
		triggerName, tableName)
}
