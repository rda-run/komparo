package diff

import (
	"fmt"
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
	var diffs []SchemaDiff

	// Tables and Columns
	diffs = append(diffs, compareTables(expected.Tables, actual.Tables)...)

	// Generic comparators
	diffs = append(diffs, compareGeneric("Index", expected.Indices, actual.Indices, 
		func(i schema.Index) string { return i.TableName + "." + i.Name },
		func(i schema.Index) string { return i.Definition })...)

	diffs = append(diffs, compareGeneric("Constraint", expected.Constraints, actual.Constraints, 
		func(c schema.Constraint) string { return c.TableName + "." + c.Name },
		func(c schema.Constraint) string { return c.Type + ": " + c.Definition })...)

	diffs = append(diffs, compareGeneric("Sequence", expected.Sequences, actual.Sequences, 
		func(s schema.Sequence) string { return s.Name },
		func(s schema.Sequence) string { return fmt.Sprintf("start=%d, min=%d, inc=%d, cycle=%v", s.Start, s.Min, s.Increment, s.Cycle) })...)

	diffs = append(diffs, compareGeneric("View", expected.Views, actual.Views, 
		func(v schema.View) string { return v.Name },
		func(v schema.View) string { return v.Definition })...)

	diffs = append(diffs, compareGeneric("MaterializedView", expected.MatViews, actual.MatViews, 
		func(v schema.MatView) string { return v.Name },
		func(v schema.MatView) string { return v.Definition })...)

	diffs = append(diffs, compareGeneric("Trigger", expected.Triggers, actual.Triggers, 
		func(t schema.Trigger) string { return t.TableName + "." + t.Name },
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

func compareTables(exp, act map[string]schema.Table) []SchemaDiff {
	var diffs []SchemaDiff

	for name, expT := range exp {
		actT, exists := act[name]
		if !exists {
			diffs = append(diffs, SchemaDiff{
				Object: name, Type: "Table", Status: "Missing", Expected: "Exists",
			})
			continue
		}

		for colName, expCol := range expT.Columns {
			actCol, colExists := actT.Columns[colName]
			if !colExists {
				diffs = append(diffs, SchemaDiff{
					Object: name + "." + colName, Type: "Column", Status: "Missing", Expected: expCol.Type,
				})
				continue
			}

			if expCol.Type != actCol.Type {
				diffs = append(diffs, SchemaDiff{
					Object: name + "." + colName, Type: "Column", Status: "Mismatch", 
					Expected: expCol.Type, Actual: actCol.Type,
				})
			}
			if expCol.Nullable != actCol.Nullable {
				diffs = append(diffs, SchemaDiff{
					Object: name + "." + colName, Type: "Column", Status: "Mismatch", 
					Expected: fmt.Sprintf("nullable=%v", expCol.Nullable), Actual: fmt.Sprintf("nullable=%v", actCol.Nullable),
				})
			}
		}

		for colName, actCol := range actT.Columns {
			if _, ok := expT.Columns[colName]; !ok {
				diffs = append(diffs, SchemaDiff{
					Object: name + "." + colName, Type: "Column", Status: "Extra", Actual: actCol.Type,
				})
			}
		}
	}

	for name := range act {
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
	for k, expV := range exp {
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
	for k, actV := range act {
		if _, ok := exp[k]; !ok {
			diffs = append(diffs, SchemaDiff{
				Object: getName(actV), Type: objType, Status: "Extra", Actual: "Exists",
			})
		}
	}
	return diffs
}
