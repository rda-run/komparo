package schema

import (
	"fmt"
	"strings"
)

// SchemaSnapshot represents the entire database schema at a point in time.
// It will be serialized to JSON.
type SchemaSnapshot struct {
	Version     string                `json:"version,omitempty"`
	PostgreSQL  string                `json:"postgresql_version"`
	Tables      map[string]Table      `json:"tables"`
	Indices     map[string]Index      `json:"indices"`
	Constraints map[string]Constraint `json:"constraints"`
	Sequences   map[string]Sequence   `json:"sequences"`
	MatViews    map[string]MatView    `json:"materialized_views"`
	Views       map[string]View       `json:"views"`
	Triggers    map[string]Trigger    `json:"triggers"`
	Functions   map[string]Function   `json:"functions"`
	Enums       map[string]Enum       `json:"enums"`
	Extensions  map[string]Extension  `json:"extensions"`
}

type Table struct {
	Name    string            `json:"name"`
	Columns map[string]Column `json:"columns"`
}

type Column struct {
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	Nullable     bool    `json:"nullable"`
	Default      *string `json:"default,omitempty"`
	MaxLen       *int    `json:"max_length,omitempty"`
	NumericPrec  *int    `json:"numeric_precision,omitempty"`
	NumericScale *int    `json:"numeric_scale,omitempty"`
	Collation    *string `json:"collation,omitempty"`
}

// FullType returns the type with its length or precision modifiers,
// e.g. "character varying(255)" or "numeric(10,2)".
func (c Column) FullType() string {
	switch strings.ToLower(c.Type) {
	case "character varying", "character", "varchar", "char", "bit", "bit varying":
		if c.MaxLen != nil && *c.MaxLen > 0 {
			return fmt.Sprintf("%s(%d)", c.Type, *c.MaxLen)
		}
	case "numeric", "decimal":
		if c.NumericPrec != nil && c.NumericScale != nil {
			return fmt.Sprintf("%s(%d,%d)", c.Type, *c.NumericPrec, *c.NumericScale)
		}
		if c.NumericPrec != nil {
			return fmt.Sprintf("%s(%d)", c.Type, *c.NumericPrec)
		}
	}
	return c.Type
}

type Index struct {
	Name       string `json:"name"`
	TableName  string `json:"table_name"`
	Definition string `json:"definition"`
}

type Constraint struct {
	Name       string `json:"name"`
	Type       string `json:"type"` // pg_constraint.contype: p, f, u, c, x, ...
	TableName  string `json:"table_name"`
	Definition string `json:"definition"`
}

// Key identifies the constraint: names are only unique per table.
func (c Constraint) Key() string { return c.TableName + "." + c.Name }

type Sequence struct {
	Name      string `json:"name"`
	Start     int64  `json:"start_value"`
	Min       int64  `json:"minimum_value"`
	Max       int64  `json:"maximum_value"`
	Increment int64  `json:"increment"`
	Cycle     bool   `json:"cycle_option"`
}

type MatView struct {
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

type View struct {
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

type Trigger struct {
	Name       string `json:"name"`
	TableName  string `json:"table_name"`
	Definition string `json:"definition"`
}

// Key identifies the trigger: names are only unique per table.
func (t Trigger) Key() string { return t.TableName + "." + t.Name }

// Normalize re-keys constraints and triggers by table and name, so snapshots
// written by older versions (keyed by name only) compare correctly.
func (s *SchemaSnapshot) Normalize() {
	s.Constraints = rekey(s.Constraints, Constraint.Key)
	s.Triggers = rekey(s.Triggers, Trigger.Key)
}

func rekey[T any](m map[string]T, key func(T) string) map[string]T {
	out := make(map[string]T, len(m))
	for _, v := range m {
		out[key(v)] = v
	}
	return out
}

type Function struct {
	Name       string `json:"name"`
	Args       string `json:"arguments"`
	ReturnType string `json:"return_type"`
	Definition string `json:"definition"`
	Volatility string `json:"volatility"`
}

type Enum struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type Extension struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
