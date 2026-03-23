package db

import "github.com/rda-run/komparo/pkg/schema"

func (e *Extractor) extractViews(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT 
			table_name,
			view_definition
		FROM information_schema.views
		WHERE table_schema = 'public'
	`
	rows, err := e.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var viewName, viewDef string
		if err := rows.Scan(&viewName, &viewDef); err != nil {
			return err
		}

		snapshot.Views[viewName] = schema.View{
			Name:       viewName,
			Definition: viewDef,
		}
	}
	return rows.Err()
}

func (e *Extractor) extractMatViews(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT 
			matviewname,
			definition
		FROM pg_matviews
		WHERE schemaname = 'public'
	`
	rows, err := e.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var mvName, mvDef string
		if err := rows.Scan(&mvName, &mvDef); err != nil {
			return err
		}

		snapshot.MatViews[mvName] = schema.MatView{
			Name:       mvName,
			Definition: mvDef,
		}
	}
	return rows.Err()
}
