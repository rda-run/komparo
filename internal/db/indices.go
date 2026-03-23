package db

import "github.com/rda-run/komparo/pkg/schema"

func (e *Extractor) extractIndices(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT 
			indexname, 
			tablename, 
			indexdef
		FROM pg_indexes
		WHERE schemaname = 'public'
	`
	rows, err := e.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var idxName, tableName, idxDef string
		if err := rows.Scan(&idxName, &tableName, &idxDef); err != nil {
			return err
		}
		
		snapshot.Indices[idxName] = schema.Index{
			Name:       idxName,
			TableName:  tableName,
			Definition: idxDef,
		}
	}
	return rows.Err()
}
