package diff

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/rda-run/komparo/pkg/schema"
)

type SchemaDiff struct {
	Object   string `json:"object"`
	Type     string `json:"type"`
	Status   string `json:"status"` // "Missing", "Extra", "Mismatch"
	Expected string `json:"expected,omitempty"`
	Actual   string `json:"actual,omitempty"`
}

func Compare(expected, actual *schema.SchemaSnapshot) []SchemaDiff {
	// Never nil, so the JSON output is "[]" instead of "null" when there are no differences.
	diffs := []SchemaDiff{}

	// Tables and Columns
	diffs = append(diffs, compareTables(expected.Tables, actual.Tables)...)

	// Generic comparators
	diffs = append(diffs, compareGeneric("Index", expected.Indices, actual.Indices,
		func(i schema.Index) string { return i.TableName + "." + i.Name },
		func(i schema.Index) string { return i.Definition })...)

	diffs = append(diffs, compareGeneric("Constraint", expected.Constraints, actual.Constraints,
		schema.Constraint.Key,
		func(c schema.Constraint) string { return c.Type + ": " + c.Definition })...)

	diffs = append(diffs, compareGeneric("Sequence", expected.Sequences, actual.Sequences,
		func(s schema.Sequence) string { return s.Name },
		func(s schema.Sequence) string {
			return fmt.Sprintf("start=%d, min=%d, inc=%d, cycle=%v", s.Start, s.Min, s.Increment, s.Cycle)
		})...)

	diffs = append(diffs, compareGeneric("View", expected.Views, actual.Views,
		func(v schema.View) string { return v.Name },
		func(v schema.View) string { return v.Definition })...)

	diffs = append(diffs, compareGeneric("MaterializedView", expected.MatViews, actual.MatViews,
		func(v schema.MatView) string { return v.Name },
		func(v schema.MatView) string { return v.Definition })...)

	diffs = append(diffs, compareGeneric("Trigger", expected.Triggers, actual.Triggers,
		schema.Trigger.Key,
		func(t schema.Trigger) string { return t.Definition })...)

	diffs = append(diffs, compareGeneric("Function", expected.Functions, actual.Functions,
		func(f schema.Function) string { return f.Name + "(" + f.Args + ")" },
		func(f schema.Function) string { return fmt.Sprintf("returns %s, %s", f.ReturnType, f.Volatility) })...)

	diffs = append(diffs, compareGeneric("Enum", expected.Enums, actual.Enums,
		func(e schema.Enum) string { return e.Name },
		func(e schema.Enum) string { return strings.Join(e.Values, ",") })...)

	diffs = append(diffs, compareGeneric("Extension", expected.Extensions, actual.Extensions,
		func(e schema.Extension) string { return e.Name },
		func(e schema.Extension) string { return e.Version })...)

	return diffs
}

// sortedKeys makes the report order deterministic between runs.
func sortedKeys[T any](m map[string]T) []string {
	return slices.Sorted(maps.Keys(m))
}

func compareTables(exp, act map[string]schema.Table) []SchemaDiff {
	var diffs []SchemaDiff

	for _, name := range sortedKeys(exp) {
		expT := exp[name]
		actT, exists := act[name]
		if !exists {
			diffs = append(diffs, SchemaDiff{
				Object: name, Type: "Table", Status: "Missing", Expected: "Exists",
			})
			continue
		}

		for _, colName := range sortedKeys(expT.Columns) {
			expCol := expT.Columns[colName]
			actCol, colExists := actT.Columns[colName]
			if !colExists {
				diffs = append(diffs, SchemaDiff{
					Object: name + "." + colName, Type: "Column", Status: "Missing", Expected: expCol.FullType(),
				})
				continue
			}

			if expType, actType := expCol.FullType(), actCol.FullType(); expType != actType {
				diffs = append(diffs, SchemaDiff{
					Object: name + "." + colName, Type: "Column", Status: "Mismatch",
					Expected: expType, Actual: actType,
				})
			}
			if expCol.Nullable != actCol.Nullable {
				diffs = append(diffs, SchemaDiff{
					Object: name + "." + colName, Type: "Column", Status: "Mismatch",
					Expected: fmt.Sprintf("nullable=%v", expCol.Nullable), Actual: fmt.Sprintf("nullable=%v", actCol.Nullable),
				})
			}
		}

		for _, colName := range sortedKeys(actT.Columns) {
			if _, ok := expT.Columns[colName]; !ok {
				diffs = append(diffs, SchemaDiff{
					Object: name + "." + colName, Type: "Column", Status: "Extra", Actual: actT.Columns[colName].FullType(),
				})
			}
		}
	}

	for _, name := range sortedKeys(act) {
		if _, ok := exp[name]; !ok {
			diffs = append(diffs, SchemaDiff{
				Object: name, Type: "Table", Status: "Extra", Actual: "Exists",
			})
		}
	}

	return diffs
}

func compareGeneric[T any](objType string, exp, act map[string]T, getName func(T) string, getDef func(T) string) []SchemaDiff {
	var diffs []SchemaDiff
	for _, k := range sortedKeys(exp) {
		expV := exp[k]
		actV, exists := act[k]
		if !exists {
			diffs = append(diffs, SchemaDiff{
				Object: getName(expV), Type: objType, Status: "Missing", Expected: "Exists",
			})
			continue
		}

		if strings.TrimSpace(getDef(expV)) != strings.TrimSpace(getDef(actV)) {
			diffs = append(diffs, SchemaDiff{
				Object: getName(expV), Type: objType, Status: "Mismatch", Expected: "Match", Actual: "Differ",
			})
		}
	}
	for _, k := range sortedKeys(act) {
		if _, ok := exp[k]; !ok {
			diffs = append(diffs, SchemaDiff{
				Object: getName(act[k]), Type: objType, Status: "Extra", Actual: "Exists",
			})
		}
	}
	return diffs
}
