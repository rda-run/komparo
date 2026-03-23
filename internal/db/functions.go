package db

import "github.com/rda-run/komparo/pkg/schema"

func (e *Extractor) extractFunctions(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT 
			r.routine_name,
			COALESCE(pg_get_function_identity_arguments(p.oid), '') AS arguments,
			r.data_type AS return_type,
			r.routine_definition AS definition,
			CASE p.provolatile 
				WHEN 'i' THEN 'IMMUTABLE' 
				WHEN 's' THEN 'STABLE' 
				WHEN 'v' THEN 'VOLATILE' 
			END AS volatility
		FROM information_schema.routines r
		JOIN pg_proc p ON r.routine_name = p.proname
		JOIN pg_namespace n ON p.pronamespace = n.oid AND n.nspname = r.routine_schema
		WHERE r.routine_schema = 'public' 
		  AND r.routine_type = 'FUNCTION'
		  AND r.external_language != 'c'
	`
	rows, err := e.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name, args string
		var retType, def, vol *string
		if err := rows.Scan(&name, &args, &retType, &def, &vol); err != nil {
			return err
		}

		definition := ""
		if def != nil {
			definition = *def
		}
		
		volatility := "VOLATILE"
		if vol != nil {
			volatility = *vol
		}
		
		returnType := ""
		if retType != nil {
			returnType = *retType
		}

		key := name + "(" + args + ")"
		snapshot.Functions[key] = schema.Function{
			Name:       name,
			Args:       args,
			ReturnType: returnType,
			Definition: definition,
			Volatility: volatility,
		}
	}
	return rows.Err()
}
