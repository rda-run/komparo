package migrator

import (
	"strconv"

	"github.com/rda-run/komparo/pkg/schema"
)

func Analyze(expected, actual *schema.SchemaSnapshot) []SchemaChange {
	var changes []SchemaChange

	changes = append(changes, analyzeTables(expected.Tables, actual.Tables)...)

	changes = append(changes, analyzeObjects("Enum", expected.Enums, actual.Enums,
		func(e, a schema.Enum) ChangeType { return "" })...)

	changes = append(changes, analyzeObjects("Extension", expected.Extensions, actual.Extensions,
		func(e, a schema.Extension) ChangeType { return alterIf(e.Version != a.Version) })...)

	changes = append(changes, analyzeObjects("Sequence", expected.Sequences, actual.Sequences,
		func(e, a schema.Sequence) ChangeType { return alterIf(e != a) })...)

	changes = append(changes, analyzeObjects("Index", expected.Indices, actual.Indices,
		func(e, a schema.Index) ChangeType { return recreateIf(e.Definition != a.Definition) })...)

	changes = append(changes, analyzeObjects("Constraint", expected.Constraints, actual.Constraints,
		func(e, a schema.Constraint) ChangeType {
			return recreateIf(e.Definition != a.Definition || e.Type != a.Type)
		})...)

	changes = append(changes, analyzeObjects("View", expected.Views, actual.Views,
		func(e, a schema.View) ChangeType { return alterIf(e.Definition != a.Definition) })...)

	changes = append(changes, analyzeObjects("MaterializedView", expected.MatViews, actual.MatViews,
		func(e, a schema.MatView) ChangeType { return recreateIf(e.Definition != a.Definition) })...)

	changes = append(changes, analyzeObjects("Function", expected.Functions, actual.Functions,
		func(e, a schema.Function) ChangeType {
			return alterIf(e.Definition != a.Definition ||
				e.ReturnType != a.ReturnType ||
				e.Volatility != a.Volatility)
		})...)

	changes = append(changes, analyzeObjects("Trigger", expected.Triggers, actual.Triggers,
		func(e, a schema.Trigger) ChangeType { return recreateIf(e.Definition != a.Definition) })...)

	return changes
}

func alterIf(differ bool) ChangeType {
	if differ {
		return ChangeAlter
	}
	return ""
}

func recreateIf(differ bool) ChangeType {
	if differ {
		return ChangeRecreate
	}
	return ""
}

// analyzeObjects emits CREATE for objects only in expected, DROP for objects
// only in actual, and whatever changed returns (if not empty) for the rest.
// The map key is the object name.
func analyzeObjects[T any](objType string, expected, actual map[string]T, changed func(exp, act T) ChangeType) []SchemaChange {
	var changes []SchemaChange

	for name, exp := range expected {
		act, exists := actual[name]
		if !exists {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeCreate, ObjectType: objType, ObjectName: name, Expected: exp,
			})
		} else if ct := changed(exp, act); ct != "" {
			changes = append(changes, SchemaChange{
				ChangeType: ct, ObjectType: objType, ObjectName: name, Expected: exp, Actual: act,
			})
		}
	}

	for name := range actual {
		if _, ok := expected[name]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeDrop, ObjectType: objType, ObjectName: name, Actual: actual[name],
			})
		}
	}

	return changes
}

func analyzeTables(expected, actual map[string]schema.Table) []SchemaChange {
	var changes []SchemaChange

	for name, expTable := range expected {
		actTable, exists := actual[name]
		if !exists {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeCreate, ObjectType: "Table", ObjectName: name, Expected: expTable,
			})
			continue
		}

		for colName, expCol := range expTable.Columns {
			expChange := ColumnChange{TableName: name, ColumnName: colName, Expected: expCol}

			actCol, colExists := actTable.Columns[colName]
			if !colExists {
				changes = append(changes, SchemaChange{
					ChangeType: ChangeCreate, ObjectType: "Column", ObjectName: name + "." + colName,
					Expected: expChange,
				})
				continue
			}

			var details []PropertyChange

			if expCol.Type != actCol.Type {
				details = append(details, PropertyChange{Property: "type", OldValue: actCol.Type, NewValue: expCol.Type})
			}

			if expCol.Nullable != actCol.Nullable {
				details = append(details, PropertyChange{
					Property: "nullable",
					OldValue: strconv.FormatBool(!actCol.Nullable),
					NewValue: strconv.FormatBool(!expCol.Nullable),
				})
			}

			if expDef, actDef := deref(expCol.Default), deref(actCol.Default); (expCol.Default == nil) != (actCol.Default == nil) || expDef != actDef {
				details = append(details, PropertyChange{Property: "default", OldValue: actDef, NewValue: expDef})
			}

			if len(details) > 0 {
				changes = append(changes, SchemaChange{
					ChangeType: ChangeAlter, ObjectType: "Column", ObjectName: name + "." + colName,
					Expected: expChange,
					Actual:   ColumnChange{TableName: name, ColumnName: colName, Actual: actCol},
					Details:  details,
				})
			}
		}

		for colName := range actTable.Columns {
			if _, ok := expTable.Columns[colName]; !ok {
				changes = append(changes, SchemaChange{
					ChangeType: ChangeDrop, ObjectType: "Column", ObjectName: name + "." + colName,
					Actual: ColumnChange{TableName: name, ColumnName: colName, Actual: actTable.Columns[colName]},
				})
			}
		}
	}

	for name := range actual {
		if _, ok := expected[name]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeDrop, ObjectType: "Table", ObjectName: name, Actual: actual[name],
			})
		}
	}

	return changes
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
