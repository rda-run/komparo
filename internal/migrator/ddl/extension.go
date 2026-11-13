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

// GenerateUpdateExtensionDDL changes the version of an installed extension.
// CREATE EXTENSION IF NOT EXISTS would be a no-op here.
func GenerateUpdateExtensionDDL(e schema.Extension) string {
	return fmt.Sprintf("ALTER EXTENSION %s UPDATE TO %s;",
		q(e.Name), pq.QuoteLiteral(e.Version))
}

func GenerateDropExtensionDDL(name string) string {
	return fmt.Sprintf("DROP EXTENSION IF EXISTS %s;", q(name))
}
