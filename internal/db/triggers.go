package db

import "github.com/rda-run/komparo/pkg/schema"

func (e *Extractor) extractTriggers(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT 
			trigger_name,
			event_object_table AS table_name,
			action_statement AS definition
		FROM information_schema.triggers
		WHERE trigger_schema = 'public'
	`
	rows, err := e.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var triggerName, tableName, definition string
		if err := rows.Scan(&triggerName, &tableName, &definition); err != nil {
			return err
		}

		snapshot.Triggers[triggerName] = schema.Trigger{
			Name:       triggerName,
			TableName:  tableName,
			Definition: definition,
		}
	}
	return rows.Err()
}
