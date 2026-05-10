package ddl

import (
	"fmt"

	"github.com/rda-run/komparo/pkg/schema"
)

func GenerateCreateOrReplaceFunctionDDL(f schema.Function) string {
	volatilityClause := ""
	if f.Volatility != "" {
		volatilityClause = f.Volatility + "\n"
	}

	return fmt.Sprintf("CREATE OR REPLACE FUNCTION %s(%s)\nRETURNS %s\nLANGUAGE plpgsql\n%sAS $$\n%s\n$$;",
		f.Name, f.Args, f.ReturnType, volatilityClause, f.Definition)
}

func GenerateDropFunctionDDL(name, args string) string {
	return fmt.Sprintf("DROP FUNCTION IF EXISTS %s(%s);", name, args)
}
