package ddl

import (
	"fmt"

	"github.com/rda-run/komparo/pkg/schema"
)

func GenerateCreateViewDDL(v schema.View) string {
	return fmt.Sprintf("CREATE OR REPLACE VIEW %s AS\n%s;",
		q(v.Name), v.Definition)
}

func GenerateCreateMatViewDDL(v schema.MatView) string {
	return fmt.Sprintf("CREATE MATERIALIZED VIEW %s AS\n%s;",
		q(v.Name), v.Definition)
}

func GenerateDropViewDDL(name string) string {
	return fmt.Sprintf("DROP VIEW IF EXISTS %s;", q(name))
}

func GenerateDropMatViewDDL(name string) string {
	return fmt.Sprintf("DROP MATERIALIZED VIEW IF EXISTS %s;", q(name))
}
