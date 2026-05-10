package ddl

import (
	"fmt"

	"github.com/rda-run/komparo/pkg/schema"
)

func GenerateCreateSequenceDDL(seq schema.Sequence) string {
	return fmt.Sprintf("CREATE SEQUENCE %s\n    START WITH %d\n    MINVALUE %d\n    MAXVALUE %d\n    INCREMENT BY %d\n    %s;",
		seq.Name,
		seq.Start,
		seq.Min,
		seq.Max,
		seq.Increment,
		map[bool]string{true: "CYCLE", false: "NO CYCLE"}[seq.Cycle])
}

func GenerateAlterSequenceDDL(seq schema.Sequence) string {
	return fmt.Sprintf("ALTER SEQUENCE %s\n    MINVALUE %d\n    MAXVALUE %d\n    INCREMENT BY %d\n    %s;",
		seq.Name,
		seq.Min,
		seq.Max,
		seq.Increment,
		map[bool]string{true: "CYCLE", false: "NO CYCLE"}[seq.Cycle])
}

func GenerateDropSequenceDDL(name string) string {
	return fmt.Sprintf("DROP SEQUENCE IF EXISTS %s;", name)
}
