package db

import "github.com/rda-run/komparo/pkg/schema"

func (e *Extractor) extractExtensions(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT 
			extname AS name,
			extversion AS version
		FROM pg_extension
	`
	rows, err := e.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name, version string
		if err := rows.Scan(&name, &version); err != nil {
			return err
		}

		snapshot.Extensions[name] = schema.Extension{
			Name:    name,
			Version: version,
		}
	}
	return rows.Err()
}
