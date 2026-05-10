package migrator

import "github.com/rda-run/komparo/pkg/schema"

type ChangeType string

const (
	ChangeCreate   ChangeType = "CREATE"
	ChangeDrop     ChangeType = "DROP"
	ChangeAlter    ChangeType = "ALTER"
	ChangeRecreate ChangeType = "RECREATE"
)

type SchemaChange struct {
	ChangeType ChangeType
	ObjectType string
	ObjectName string
	Expected   interface{}
	Actual     interface{}
	Details    []PropertyChange
}

type PropertyChange struct {
	Property string
	OldValue string
	NewValue string
}

type TableChange struct {
	TableName string
	Expected  schema.Table
	Actual    schema.Table
}

type ColumnChange struct {
	TableName  string
	ColumnName string
	Expected   schema.Column
	Actual     schema.Column
}

type IndexChange struct {
	Expected schema.Index
	Actual   schema.Index
}

type ConstraintChange struct {
	Expected schema.Constraint
	Actual   schema.Constraint
}

type SequenceChange struct {
	Expected schema.Sequence
	Actual   schema.Sequence
}

type ViewChange struct {
	Expected schema.View
	Actual   schema.View
}

type MatViewChange struct {
	Expected schema.MatView
	Actual   schema.MatView
}

type TriggerChange struct {
	Expected schema.Trigger
	Actual   schema.Trigger
}

type FunctionChange struct {
	Expected schema.Function
	Actual   schema.Function
}

type EnumChange struct {
	Expected schema.Enum
	Actual   schema.Enum
}

type ExtensionChange struct {
	Expected schema.Extension
	Actual   schema.Extension
}
