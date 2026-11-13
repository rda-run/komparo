package db

import "github.com/rda-run/komparo/pkg/schema"

// Reads pg_proc directly, producing the same values information_schema.routines
// would for the function owner. The previous query joined routines to pg_proc
// by name only (mixing up overloads) and got a NULL body for non-owners.
func (e *Extractor) extractFunctions(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT
			p.proname,
			COALESCE(pg_get_function_identity_arguments(p.oid), '') AS arguments,
			CASE
				WHEN t.typelem <> 0 AND t.typlen = -1 THEN 'ARRAY'
				WHEN tn.nspname = 'pg_catalog' THEN format_type(t.oid, NULL)
				ELSE 'USER-DEFINED'
			END AS return_type,
			p.prosrc AS definition,
			CASE p.provolatile
				WHEN 'i' THEN 'IMMUTABLE'
				WHEN 's' THEN 'STABLE'
				WHEN 'v' THEN 'VOLATILE'
			END AS volatility
		FROM pg_proc p
		JOIN pg_namespace n ON n.oid = p.pronamespace
		JOIN pg_type t ON t.oid = p.prorettype
		JOIN pg_namespace tn ON tn.oid = t.typnamespace
		WHERE n.nspname = 'public'
		  AND p.prokind = 'f'
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
