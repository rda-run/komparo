package ddl

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/rda-run/komparo/pkg/schema"
)

func GenerateCreateExtensionDDL(e schema.Extension) string {
	versionClause := ""
	if e.Version != "" {
		versionClause = " WITH VERSION " + pq.QuoteLiteral(e.Version)
	}
	return fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s%s;",
		q(e.Name), versionClause)
}

func GenerateDropExtensionDDL(name string) string {
	return fmt.Sprintf("DROP EXTENSION IF EXISTS %s;", q(name))
}
