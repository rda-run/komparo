package migrator

import (
	"strings"
	"testing"

	"github.com/rda-run/komparo/pkg/schema"
)

func intp(i int) *int { return &i }

func TestFixScript(t *testing.T) {
	expected := &schema.SchemaSnapshot{
		Tables: map[string]schema.Table{
			"order": {Name: "order", Columns: map[string]schema.Column{
				"id":    {Name: "id", Type: "integer"},
				"email": {Name: "email", Type: "character varying", MaxLen: intp(255)},
			}},
			"users": {Name: "users", Columns: map[string]schema.Column{
				"email": {Name: "email", Type: "character varying", MaxLen: intp(255)},
			}},
		},
		Indices: map[string]schema.Index{
			"order_pkey": {Name: "order_pkey", TableName: "order", Definition: `CREATE UNIQUE INDEX order_pkey ON public."order" USING btree (id)`},
			"idx_email":  {Name: "idx_email", TableName: "order", Definition: `CREATE INDEX idx_email ON public."order" USING btree (email)`},
		},
		Constraints: map[string]schema.Constraint{
			"order.fk":         {Name: "fk", Type: "f", TableName: "order", Definition: "FOREIGN KEY (id) REFERENCES users(id)"},
			"order.order_pkey": {Name: "order_pkey", Type: "p", TableName: "order", Definition: "PRIMARY KEY (id)"},
		},
		Enums:      map[string]schema.Enum{"status": {Name: "status", Values: []string{"new", "it's", "done"}}},
		Extensions: map[string]schema.Extension{"uuid-ossp": {Name: "uuid-ossp", Version: "1.1"}},
	}
	actual := &schema.SchemaSnapshot{
		Tables: map[string]schema.Table{
			"users": {Name: "users", Columns: map[string]schema.Column{
				"email": {Name: "email", Type: "character varying", MaxLen: intp(50)},
			}},
		},
		Enums:      map[string]schema.Enum{"status": {Name: "status", Values: []string{"new", "old"}}},
		Extensions: map[string]schema.Extension{"uuid-ossp": {Name: "uuid-ossp", Version: "1.0"}},
	}

	sql := GenerateFixScript(OrderChanges(Analyze(expected, actual)), "postgresql://u:secret@h/db", "snap.json")

	for _, want := range []string{
		"\nBEGIN;\n",
		`CREATE TABLE "order" (`,
		`"email" character varying(255) NOT NULL`,
		`ALTER TABLE "users" ALTER COLUMN "email" TYPE character varying(255);`,
		`ALTER TYPE "status" ADD VALUE IF NOT EXISTS 'it''s' AFTER 'new';`,
		`ALTER TYPE "status" ADD VALUE IF NOT EXISTS 'done' AFTER 'it''s';`,
		"values not in the expected schema (old)",
		`ALTER EXTENSION "uuid-ossp" UPDATE TO '1.1';`,
		"CREATE INDEX idx_email",
	} {
		if !strings.Contains(sql, want) {
			t.Errorf("script is missing %q", want)
		}
	}

	for _, unwanted := range []string{
		"secret",                         // password must be redacted
		"CREATE UNIQUE INDEX order_pkey", // created by the PRIMARY KEY constraint
	} {
		if strings.Contains(sql, unwanted) {
			t.Errorf("script must not contain %q", unwanted)
		}
	}

	if pk, fk := strings.Index(sql, "PRIMARY KEY"), strings.Index(sql, "FOREIGN KEY"); pk < 0 || fk < pk {
		t.Errorf("foreign key must come after the primary key (pk=%d fk=%d)", pk, fk)
	}

	if again := GenerateFixScript(OrderChanges(Analyze(expected, actual)), "postgresql://u:secret@h/db", "snap.json"); stripHeader(again) != stripHeader(sql) {
		t.Error("script is not deterministic")
	}

	if t.Failed() {
		t.Log(sql)
	}
}

// stripHeader drops the header, which carries a timestamp.
func stripHeader(sql string) string {
	return sql[strings.Index(sql, "BEGIN;"):]
}
