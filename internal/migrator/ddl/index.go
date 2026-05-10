package ddl

import (
	"fmt"

	"github.com/rda-run/komparo/pkg/schema"
)

func GenerateCreateIndexDDL(idx schema.Index) string {
	return idx.Definition + ";"
}

func GenerateDropIndexDDL(idxName string) string {
	return fmt.Sprintf("DROP INDEX IF EXISTS %s;", idxName)
}
