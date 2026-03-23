package db

import (
	"database/sql"
	"github.com/rda-run/komparo/pkg/schema"
)

func (e *Extractor) extractSequences(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT 
			sequence_name,
			start_value::bigint,
			minimum_value::bigint,
			maximum_value::bigint,
			increment::bigint,
			cycle_option
		FROM information_schema.sequences
		WHERE sequence_schema = 'public'
	`
	rows, err := e.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var seqName string
		var start, min, max, inc sql.NullInt64
		var cycleStr string

		if err := rows.Scan(&seqName, &start, &min, &max, &inc, &cycleStr); err != nil {
			return err
		}

		snapshot.Sequences[seqName] = schema.Sequence{
			Name:      seqName,
			Start:     start.Int64,
			Min:       min.Int64,
			Max:       max.Int64,
			Increment: inc.Int64,
			Cycle:     cycleStr == "YES",
		}
	}
	return rows.Err()
}
