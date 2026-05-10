package ddl

import (
	"fmt"
	"strings"

	"github.com/rda-run/komparo/pkg/schema"
)

func GenerateCreateEnumDDL(e schema.Enum) string {
	quoted := make([]string, len(e.Values))
	for i, v := range e.Values {
		quoted[i] = fmt.Sprintf("'%s'", v)
	}
	return fmt.Sprintf("CREATE TYPE %s AS ENUM (%s);",
		e.Name, strings.Join(quoted, ", "))
}

func GenerateDropEnumDDL(name string) string {
	return fmt.Sprintf("DROP TYPE IF EXISTS %s;", name)
}

func GenerateAddEnumValueDDL(enumName, value string) string {
	return fmt.Sprintf("ALTER TYPE %s ADD VALUE '%s';", enumName, value)
}
