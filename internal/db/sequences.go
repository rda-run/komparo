package db

import "github.com/rda-run/komparo/pkg/schema"

// Reads pg_catalog instead of information_schema.sequences: the latter hides
// sequences on which the user has no privilege (typical for read-only users).
func (e *Extractor) extractSequences(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT
			c.relname,
			s.seqstart,
			s.seqmin,
			s.seqmax,
			s.seqincrement,
			s.seqcycle
		FROM pg_sequence s
		JOIN pg_class c ON c.oid = s.seqrelid
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = 'public'
		  AND NOT EXISTS (
			SELECT 1 FROM pg_depend d
			WHERE d.classid = 'pg_class'::regclass AND d.objid = c.oid AND d.deptype = 'i'
		  ) -- identity column sequences are skipped, as in information_schema.sequences
	`
	rows, err := e.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var seqName string
		var start, min, max, inc int64
		var cycle bool

		if err := rows.Scan(&seqName, &start, &min, &max, &inc, &cycle); err != nil {
			return err
		}

		snapshot.Sequences[seqName] = schema.Sequence{
			Name:      seqName,
			Start:     start,
			Min:       min,
			Max:       max,
			Increment: inc,
			Cycle:     cycle,
		}
	}
	return rows.Err()
}
