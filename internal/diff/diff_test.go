package diff

import (
	"github.com/rda-run/komparo/pkg/schema"
	"testing"
)

func TestCompareTables(t *testing.T) {
	expected := &schema.SchemaSnapshot{
		Tables: map[string]schema.Table{
			"users": {
				Name: "users",
				Columns: map[string]schema.Column{
					"id":  {Name: "id", Type: "integer", Nullable: false},
					"age": {Name: "age", Type: "integer", Nullable: true}, // Missing in actual
				},
			},
		},
	}

	actual := &schema.SchemaSnapshot{
		Tables: map[string]schema.Table{
			"users": {
				Name: "users",
				Columns: map[string]schema.Column{
					"id":   {Name: "id", Type: "bigint", Nullable: false}, // Mismatch
					"name": {Name: "name", Type: "text", Nullable: true},  // Extra
				},
			},
		},
	}

	diffs := Compare(expected, actual)

	if len(diffs) != 3 {
		t.Fatalf("expected 3 diffs, got %d", len(diffs))
	}

	missingFound, extraFound, mismatchFound := false, false, false
	for _, d := range diffs {
		if d.Status == "Missing" && d.Object == "users.age" {
			missingFound = true
		}
		if d.Status == "Extra" && d.Object == "users.name" {
			extraFound = true
		}
		if d.Status == "Mismatch" && d.Object == "users.id" {
			mismatchFound = true
		}
	}

	if !missingFound || !extraFound || !mismatchFound {
		t.Errorf("Diff logic failed to correctly categorize table modifications")
	}
}

func TestCompareDetectsLengthAndCollidingNames(t *testing.T) {
	n := func(i int) *int { return &i }
	expected := &schema.SchemaSnapshot{
		Tables: map[string]schema.Table{"users": {Name: "users", Columns: map[string]schema.Column{
			"email": {Name: "email", Type: "character varying", MaxLen: n(255)},
			"total": {Name: "total", Type: "numeric", NumericPrec: n(10), NumericScale: n(2)},
			"id":    {Name: "id", Type: "integer", NumericPrec: n(32), NumericScale: n(0)},
		}}},
		// Snapshot written by an older version: keyed by name only.
		Triggers: map[string]schema.Trigger{
			"set_updated_at": {Name: "set_updated_at", TableName: "users", Definition: "EXECUTE FUNCTION f()"},
		},
	}
	expected.Normalize()

	actual := &schema.SchemaSnapshot{
		Tables: map[string]schema.Table{"users": {Name: "users", Columns: map[string]schema.Column{
			"email": {Name: "email", Type: "character varying", MaxLen: n(50)},
			"total": {Name: "total", Type: "numeric", NumericPrec: n(12), NumericScale: n(4)},
			"id":    {Name: "id", Type: "integer", NumericPrec: n(32), NumericScale: n(0)},
		}}},
		Triggers: map[string]schema.Trigger{
			"users.set_updated_at": {Name: "set_updated_at", TableName: "users", Definition: "EXECUTE FUNCTION f()"},
			"order.set_updated_at": {Name: "set_updated_at", TableName: "order", Definition: "EXECUTE FUNCTION f()"},
		},
	}

	want := []SchemaDiff{
		{Object: "users.email", Type: "Column", Status: "Mismatch", Expected: "character varying(255)", Actual: "character varying(50)"},
		{Object: "users.total", Type: "Column", Status: "Mismatch", Expected: "numeric(10,2)", Actual: "numeric(12,4)"},
		{Object: "order.set_updated_at", Type: "Trigger", Status: "Extra", Actual: "Exists"},
	}
	got := Compare(expected, actual)
	if len(got) != len(want) {
		t.Fatalf("got %+v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("diff %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}
