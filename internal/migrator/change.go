package migrator

import "github.com/rda-run/komparo/pkg/schema"

type ChangeType string

const (
	ChangeCreate   ChangeType = "CREATE"
	ChangeDrop     ChangeType = "DROP"
	ChangeAlter    ChangeType = "ALTER"
	ChangeRecreate ChangeType = "RECREATE"
)

// SchemaChange describes one difference. Expected and Actual hold the schema
// object (schema.Table, schema.Index, ...), or a ColumnChange for columns.
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

type ColumnChange struct {
	TableName  string
	ColumnName string
	Expected   schema.Column
	Actual     schema.Column
}
