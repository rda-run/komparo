package migrator

import (
	"github.com/rda-run/komparo/pkg/schema"
)

func Analyze(expected, actual *schema.SchemaSnapshot) []SchemaChange {
	var changes []SchemaChange

	changes = append(changes, analyzeTables(expected.Tables, actual.Tables)...)
	changes = append(changes, analyzeEnums(expected.Enums, actual.Enums)...)
	changes = append(changes, analyzeExtensions(expected.Extensions, actual.Extensions)...)
	changes = append(changes, analyzeSequences(expected.Sequences, actual.Sequences)...)
	changes = append(changes, analyzeIndices(expected.Indices, actual.Indices)...)
	changes = append(changes, analyzeConstraints(expected.Constraints, actual.Constraints)...)
	changes = append(changes, analyzeViews(expected.Views, actual.Views)...)
	changes = append(changes, analyzeMatViews(expected.MatViews, actual.MatViews)...)
	changes = append(changes, analyzeFunctions(expected.Functions, actual.Functions)...)
	changes = append(changes, analyzeTriggers(expected.Triggers, actual.Triggers)...)

	return changes
}

func analyzeTables(expected, actual map[string]schema.Table) []SchemaChange {
	var changes []SchemaChange

	for name, expTable := range expected {
		actTable, exists := actual[name]
		if !exists {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeCreate,
				ObjectType: "Table",
				ObjectName: name,
				Expected:   expTable,
				Actual:     nil,
			})
			continue
		}

		for colName, expCol := range expTable.Columns {
			actCol, colExists := actTable.Columns[colName]
			if !colExists {
				changes = append(changes, SchemaChange{
					ChangeType: ChangeCreate,
					ObjectType: "Column",
					ObjectName:   name + "." + colName,
					Expected: ColumnChange{
						TableName:  name,
						ColumnName: colName,
						Expected:   expCol,
					},
					Actual: nil,
				})
				continue
			}

			var details []PropertyChange

			if expCol.Type != actCol.Type {
				details = append(details, PropertyChange{
					Property: "type",
					OldValue: actCol.Type,
					NewValue: expCol.Type,
				})
			}

			if expCol.Nullable != actCol.Nullable {
				details = append(details, PropertyChange{
					Property: "nullable",
					OldValue: boolToString(!actCol.Nullable),
					NewValue: boolToString(!expCol.Nullable),
				})
			}

			if (expCol.Default == nil && actCol.Default != nil) ||
				(expCol.Default != nil && actCol.Default == nil) ||
				(expCol.Default != nil && actCol.Default != nil && *expCol.Default != *actCol.Default) {
				oldVal := ""
				if actCol.Default != nil {
					oldVal = *actCol.Default
				}
				newVal := ""
				if expCol.Default != nil {
					newVal = *expCol.Default
				}
				details = append(details, PropertyChange{
					Property: "default",
					OldValue: oldVal,
					NewValue: newVal,
				})
			}

			if len(details) > 0 {
				changes = append(changes, SchemaChange{
					ChangeType: ChangeAlter,
					ObjectType: "Column",
					ObjectName: name + "." + colName,
					Expected: ColumnChange{
						TableName:  name,
						ColumnName: colName,
						Expected:   expCol,
					},
					Actual: ColumnChange{
						TableName:  name,
						ColumnName: colName,
						Actual:     actCol,
					},
					Details: details,
				})
			}
		}

		for colName := range actTable.Columns {
			if _, ok := expTable.Columns[colName]; !ok {
				changes = append(changes, SchemaChange{
					ChangeType: ChangeDrop,
					ObjectType: "Column",
					ObjectName: name + "." + colName,
					Expected:   nil,
					Actual: ColumnChange{
						TableName:  name,
						ColumnName: colName,
						Actual:     actTable.Columns[colName],
					},
				})
			}
		}
	}

	for name, actTable := range actual {
		if _, ok := expected[name]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeDrop,
				ObjectType: "Table",
				ObjectName: name,
				Expected:   nil,
				Actual:     actTable,
			})
		}
	}

	return changes
}

func analyzeEnums(expected, actual map[string]schema.Enum) []SchemaChange {
	var changes []SchemaChange

	for name, expEnum := range expected {
		if _, ok := actual[name]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeCreate,
				ObjectType: "Enum",
				ObjectName: name,
				Expected:   expEnum,
				Actual:     nil,
			})
		}
	}

	for name, actEnum := range actual {
		if _, ok := expected[name]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeDrop,
				ObjectType: "Enum",
				ObjectName: name,
				Expected:   nil,
				Actual:     actEnum,
			})
		}
	}

	return changes
}

func analyzeExtensions(expected, actual map[string]schema.Extension) []SchemaChange {
	var changes []SchemaChange

	for name, expExt := range expected {
		actExt, exists := actual[name]
		if !exists {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeCreate,
				ObjectType: "Extension",
				ObjectName: name,
				Expected:   expExt,
				Actual:     nil,
			})
		} else if expExt.Version != actExt.Version {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeAlter,
				ObjectType: "Extension",
				ObjectName: name,
				Expected:   expExt,
				Actual:     actExt,
				Details: []PropertyChange{{
					Property: "version",
					OldValue: actExt.Version,
					NewValue: expExt.Version,
				}},
			})
		}
	}

	for name, actExt := range actual {
		if _, ok := expected[name]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeDrop,
				ObjectType: "Extension",
				ObjectName: name,
				Expected:   nil,
				Actual:     actExt,
			})
		}
	}

	return changes
}

func analyzeSequences(expected, actual map[string]schema.Sequence) []SchemaChange {
	var changes []SchemaChange

	for name, expSeq := range expected {
		actSeq, exists := actual[name]
		if !exists {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeCreate,
				ObjectType: "Sequence",
				ObjectName: name,
				Expected:   expSeq,
				Actual:     nil,
			})
		} else if sequencesDiffer(expSeq, actSeq) {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeAlter,
				ObjectType: "Sequence",
				ObjectName: name,
				Expected:   expSeq,
				Actual:     actSeq,
			})
		}
	}

	for name, actSeq := range actual {
		if _, ok := expected[name]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeDrop,
				ObjectType: "Sequence",
				ObjectName: name,
				Expected:   nil,
				Actual:     actSeq,
			})
		}
	}

	return changes
}

func sequencesDiffer(a, b schema.Sequence) bool {
	return a.Start != b.Start || a.Min != b.Min || a.Max != b.Max ||
		a.Increment != b.Increment || a.Cycle != b.Cycle
}

func analyzeIndices(expected, actual map[string]schema.Index) []SchemaChange {
	var changes []SchemaChange

	for name, expIdx := range expected {
		actIdx, exists := actual[name]
		if !exists {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeCreate,
				ObjectType: "Index",
				ObjectName: expIdx.TableName + "." + name,
				Expected:   expIdx,
				Actual:     nil,
			})
		} else if expIdx.Definition != actIdx.Definition {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeRecreate,
				ObjectType: "Index",
				ObjectName: expIdx.TableName + "." + name,
				Expected:   expIdx,
				Actual:     actIdx,
			})
		}
	}

	for name, actIdx := range actual {
		if _, ok := expected[name]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeDrop,
				ObjectType: "Index",
				ObjectName: actIdx.TableName + "." + name,
				Expected:   nil,
				Actual:     actIdx,
			})
		}
	}

	return changes
}

func analyzeConstraints(expected, actual map[string]schema.Constraint) []SchemaChange {
	var changes []SchemaChange

	for name, expCon := range expected {
		actCon, exists := actual[name]
		if !exists {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeCreate,
				ObjectType: "Constraint",
				ObjectName: expCon.TableName + "." + name,
				Expected:   expCon,
				Actual:     nil,
			})
		} else if expCon.Definition != actCon.Definition || expCon.Type != actCon.Type {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeRecreate,
				ObjectType: "Constraint",
				ObjectName: expCon.TableName + "." + name,
				Expected:   expCon,
				Actual:     actCon,
			})
		}
	}

	for name, actCon := range actual {
		if _, ok := expected[name]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeDrop,
				ObjectType: "Constraint",
				ObjectName: actCon.TableName + "." + name,
				Expected:   nil,
				Actual:     actCon,
			})
		}
	}

	return changes
}

func analyzeViews(expected, actual map[string]schema.View) []SchemaChange {
	var changes []SchemaChange

	for name, expView := range expected {
		actView, exists := actual[name]
		if !exists {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeCreate,
				ObjectType: "View",
				ObjectName: name,
				Expected:   expView,
				Actual:     nil,
			})
		} else if expView.Definition != actView.Definition {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeAlter,
				ObjectType: "View",
				ObjectName: name,
				Expected:   expView,
				Actual:     actView,
			})
		}
	}

	for name, actView := range actual {
		if _, ok := expected[name]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeDrop,
				ObjectType: "View",
				ObjectName: name,
				Expected:   nil,
				Actual:     actView,
			})
		}
	}

	return changes
}

func analyzeMatViews(expected, actual map[string]schema.MatView) []SchemaChange {
	var changes []SchemaChange

	for name, expView := range expected {
		actView, exists := actual[name]
		if !exists {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeCreate,
				ObjectType: "MaterializedView",
				ObjectName: name,
				Expected:   expView,
				Actual:     nil,
			})
		} else if expView.Definition != actView.Definition {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeRecreate,
				ObjectType: "MaterializedView",
				ObjectName: name,
				Expected:   expView,
				Actual:     actView,
			})
		}
	}

	for name, actView := range actual {
		if _, ok := expected[name]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeDrop,
				ObjectType: "MaterializedView",
				ObjectName: name,
				Expected:   nil,
				Actual:     actView,
			})
		}
	}

	return changes
}

func analyzeFunctions(expected, actual map[string]schema.Function) []SchemaChange {
	var changes []SchemaChange

	key := func(f schema.Function) string {
		return f.Name + "(" + f.Args + ")"
	}

	for _, expFunc := range expected {
		k := key(expFunc)
		actFunc, exists := actual[k]
		if !exists {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeCreate,
				ObjectType: "Function",
				ObjectName: k,
				Expected:   expFunc,
				Actual:     nil,
			})
		} else if expFunc.Definition != actFunc.Definition ||
			expFunc.ReturnType != actFunc.ReturnType ||
			expFunc.Volatility != actFunc.Volatility {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeAlter,
				ObjectType: "Function",
				ObjectName: k,
				Expected:   expFunc,
				Actual:     actFunc,
			})
		}
	}

	for _, actFunc := range actual {
		k := key(actFunc)
		if _, ok := expected[k]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeDrop,
				ObjectType: "Function",
				ObjectName: k,
				Expected:   nil,
				Actual:     actFunc,
			})
		}
	}

	return changes
}

func analyzeTriggers(expected, actual map[string]schema.Trigger) []SchemaChange {
	var changes []SchemaChange

	for name, expTrig := range expected {
		actTrig, exists := actual[name]
		if !exists {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeCreate,
				ObjectType: "Trigger",
				ObjectName: expTrig.TableName + "." + name,
				Expected:   expTrig,
				Actual:     nil,
			})
		} else if expTrig.Definition != actTrig.Definition {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeRecreate,
				ObjectType: "Trigger",
				ObjectName: expTrig.TableName + "." + name,
				Expected:   expTrig,
				Actual:     actTrig,
			})
		}
	}

	for name, actTrig := range actual {
		if _, ok := expected[name]; !ok {
			changes = append(changes, SchemaChange{
				ChangeType: ChangeDrop,
				ObjectType: "Trigger",
				ObjectName: actTrig.TableName + "." + name,
				Expected:   nil,
				Actual:     actTrig,
			})
		}
	}

	return changes
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
