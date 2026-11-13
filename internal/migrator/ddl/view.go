package ddl

import (
	"fmt"
	"strings"

	"github.com/rda-run/komparo/pkg/schema"
)

// pg_get_viewdef already ends the query with ";".
func body(definition string) string {
	return strings.TrimSuffix(strings.TrimSpace(definition), ";")
}

func GenerateCreateViewDDL(v schema.View) string {
	return fmt.Sprintf("CREATE OR REPLACE VIEW %s AS\n%s;",
		q(v.Name), body(v.Definition))
}

func GenerateCreateMatViewDDL(v schema.MatView) string {
	return fmt.Sprintf("CREATE MATERIALIZED VIEW %s AS\n%s;",
		q(v.Name), body(v.Definition))
}

func GenerateDropViewDDL(name string) string {
	return fmt.Sprintf("DROP VIEW IF EXISTS %s;", q(name))
}

func GenerateDropMatViewDDL(name string) string {
	return fmt.Sprintf("DROP MATERIALIZED VIEW IF EXISTS %s;", q(name))
}
