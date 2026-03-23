package db

import "github.com/rda-run/komparo/pkg/schema"

func (e *Extractor) extractEnums(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT 
			t.typname AS enum_name,
			e.enumlabel AS enum_value
		FROM pg_type t
		JOIN pg_enum e ON t.oid = e.enumtypid
		JOIN pg_namespace n ON n.oid = t.typnamespace
		WHERE n.nspname = 'public'
		ORDER BY t.typname, e.enumsortorder
	`
	rows, err := e.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	enumMap := make(map[string][]string)

	for rows.Next() {
		var enumName, enumValue string
		if err := rows.Scan(&enumName, &enumValue); err != nil {
			return err
		}
		enumMap[enumName] = append(enumMap[enumName], enumValue)
	}

	for name, values := range enumMap {
		snapshot.Enums[name] = schema.Enum{
			Name:   name,
			Values: values,
		}
	}
	return rows.Err()
}
