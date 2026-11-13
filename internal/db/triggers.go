package db

import "github.com/rda-run/komparo/pkg/schema"

// Reads pg_catalog instead of information_schema.triggers: the latter only
// shows triggers on tables where the user has a privilege other than SELECT,
// so read-only users would see none.
func (e *Extractor) extractTriggers(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT
			t.tgname,
			c.relname,
			substring(pg_get_triggerdef(t.oid) from 'EXECUTE (?:FUNCTION|PROCEDURE) .*$')
		FROM pg_trigger t
		JOIN pg_class c ON c.oid = t.tgrelid
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = 'public'
		  AND NOT t.tgisinternal
		  AND (t.tgtype & 28) <> 0 -- INSERT, DELETE or UPDATE (same set as information_schema.triggers)
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
