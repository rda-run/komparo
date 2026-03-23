package diff

import (
	"testing"
	"github.com/rda-run/komparo/pkg/schema"
)

func TestCompareTables(t *testing.T) {
	expected := &schema.SchemaSnapshot{
		Tables: map[string]schema.Table{
			"users": {
				Name: "users",
				Columns: map[string]schema.Column{
					"id":   {Name: "id", Type: "integer", Nullable: false},
					"age":  {Name: "age", Type: "integer", Nullable: true}, // Missing in actual
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
					"name": {Name: "name", Type: "text", Nullable: true}, // Extra
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
