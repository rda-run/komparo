package ddl

import (
	"fmt"

	"github.com/rda-run/komparo/pkg/schema"
)

func GenerateCreateExtensionDDL(e schema.Extension) string {
	versionClause := ""
	if e.Version != "" {
		versionClause = fmt.Sprintf(" WITH VERSION '%s'", e.Version)
	}
	return fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s%s;",
		e.Name, versionClause)
}

func GenerateDropExtensionDDL(name string) string {
	return fmt.Sprintf("DROP EXTENSION IF EXISTS %s;", name)
}
