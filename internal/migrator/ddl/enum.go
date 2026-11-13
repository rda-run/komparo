package ddl

import (
	"fmt"
	"slices"
	"strings"

	"github.com/lib/pq"
	"github.com/rda-run/komparo/pkg/schema"
)

func GenerateCreateEnumDDL(e schema.Enum) string {
	quoted := make([]string, len(e.Values))
	for i, v := range e.Values {
		quoted[i] = pq.QuoteLiteral(v)
	}
	return fmt.Sprintf("CREATE TYPE %s AS ENUM (%s);",
		q(e.Name), strings.Join(quoted, ", "))
}

func GenerateDropEnumDDL(name string) string {
	return fmt.Sprintf("DROP TYPE IF EXISTS %s;", q(name))
}

// GenerateAlterEnumDDL adds the values missing from actual, in the expected
// position. PostgreSQL cannot remove or reorder enum values, so those cases
// only produce a warning comment.
func GenerateAlterEnumDDL(expected, actual schema.Enum) string {
	var lines []string
	current := slices.Clone(actual.Values)

	for i, v := range expected.Values {
		if slices.Contains(current, v) {
			continue
		}
		position := ""
		if i > 0 {
			position = " AFTER " + pq.QuoteLiteral(expected.Values[i-1])
		} else if len(current) > 0 {
			position = " BEFORE " + pq.QuoteLiteral(current[0])
		}
		lines = append(lines, fmt.Sprintf("ALTER TYPE %s ADD VALUE IF NOT EXISTS %s%s;",
			q(expected.Name), pq.QuoteLiteral(v), position))
		current = append(current, v)
	}
	if len(lines) > 0 {
		lines = append([]string{"-- NOTE: a new enum value cannot be used until this transaction is committed."}, lines...)
	}

	var extra []string
	for _, v := range actual.Values {
		if !slices.Contains(expected.Values, v) {
			extra = append(extra, v)
		}
	}
	if len(extra) > 0 {
		lines = append(lines, fmt.Sprintf("-- WARNING: enum %s has values not in the expected schema (%s). PostgreSQL cannot drop enum values; recreate the type manually.",
			q(expected.Name), strings.Join(extra, ", ")))
	} else if len(lines) == 0 {
		lines = append(lines, fmt.Sprintf("-- WARNING: enum %s has the expected values in a different order. PostgreSQL cannot reorder enum values; recreate the type manually.",
			q(expected.Name)))
	}

	return strings.Join(lines, "\n")
}
