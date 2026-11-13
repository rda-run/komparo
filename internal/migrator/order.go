package migrator

import (
	"cmp"
	"slices"
)

var sectionOrder = map[ChangeType]int{ChangeDrop: 0, ChangeAlter: 1, ChangeRecreate: 2, ChangeCreate: 3}

// Creation order; drops run in reverse.
var typeOrder = map[string]int{
	"Extension": 0, "Enum": 1, "Sequence": 2, "Table": 3, "Column": 4, "Constraint": 5,
	"Index": 6, "View": 7, "MaterializedView": 8, "Function": 9, "Trigger": 10,
}

// OrderChanges sorts by section (DROP, ALTER, RECREATE, CREATE) and then by
// dependency order of the object type. The sort is stable, so the name order
// produced by Analyze is kept within each group.
func OrderChanges(changes []SchemaChange) []SchemaChange {
	result := slices.Clone(changes)
	slices.SortStableFunc(result, func(a, b SchemaChange) int {
		return cmp.Or(
			cmp.Compare(sectionOrder[a.ChangeType], sectionOrder[b.ChangeType]),
			cmp.Compare(rank(a), rank(b)),
		)
	})
	return result
}

func rank(c SchemaChange) int {
	r := typeOrder[c.ObjectType]
	if c.ChangeType == ChangeDrop {
		return -r
	}
	return r
}
