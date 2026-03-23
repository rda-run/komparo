package schema

// SchemaSnapshot represents the entire database schema at a point in time.
// It will be serialized to JSON.
type SchemaSnapshot struct {
	Version      string                 `json:"version,omitempty"`
	PostgreSQL   string                 `json:"postgresql_version"`
	Tables       map[string]Table       `json:"tables"`
	Indices      map[string]Index       `json:"indices"`
	Constraints  map[string]Constraint  `json:"constraints"`
	Sequences    map[string]Sequence    `json:"sequences"`
	MatViews     map[string]MatView     `json:"materialized_views"`
	Views        map[string]View        `json:"views"`
	Triggers     map[string]Trigger     `json:"triggers"`
	Functions    map[string]Function    `json:"functions"`
	Enums        map[string]Enum        `json:"enums"`
	Extensions   map[string]Extension   `json:"extensions"`
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

type Index struct {
	Name       string `json:"name"`
	TableName  string `json:"table_name"`
	Definition string `json:"definition"`
}

type Constraint struct {
	Name       string `json:"name"`
	Type       string `json:"type"` // PRIMARY KEY, FOREIGN KEY, UNIQUE, CHECK, EXCLUDE
	TableName  string `json:"table_name"`
	Definition string `json:"definition"`
}

type Sequence struct {
	Name      string `json:"name"`
	Start     int64  `json:"start_value"`
	Min       int64  `json:"minimum_value"`
	Max       int64  `json:"maximum_value"`
	Increment int64  `json:"increment"`
	Cycle     bool   `json:"cycle_option"`
}

type MatView struct {
	Name       string   `json:"name"`
	Definition string   `json:"definition"`
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
