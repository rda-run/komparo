package db

import "github.com/rda-run/komparo/pkg/schema"

func (e *Extractor) extractConstraints(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT 
			conname AS constraint_name,
			contype AS constraint_type,
			cl.relname AS table_name,
			pg_get_constraintdef(c.oid) AS definition
		FROM pg_constraint c
		JOIN pg_class cl ON c.conrelid = cl.oid
		JOIN pg_namespace n ON cl.relnamespace = n.oid
		WHERE n.nspname = 'public'
	`
	rows, err := e.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var conName, conType, tableName, conDef string
		if err := rows.Scan(&conName, &conType, &tableName, &conDef); err != nil {
			return err
		}

		snapshot.Constraints[conName] = schema.Constraint{
			Name:       conName,
			Type:       conType,
			TableName:  tableName,
			Definition: conDef,
		}
	}
	return rows.Err()
}
