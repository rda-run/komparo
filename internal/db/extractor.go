package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rda-run/komparo/pkg/schema"
)

// Extractor handles retrieving schema information from PostgreSQL
type Extractor struct {
	db *sql.DB
}

// NewExtractor creates a new Extractor connected to the given database URL
func NewExtractor(dbURL string) (*Extractor, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Extractor{db: db}, nil
}

// Close closes the underlying database connection
func (e *Extractor) Close() error {
	return e.db.Close()
}

// ExtractFullSchema returns a complete SchemaSnapshot
func (e *Extractor) ExtractFullSchema() (*schema.SchemaSnapshot, error) {
	snapshot := &schema.SchemaSnapshot{
		Tables:      make(map[string]schema.Table),
		Indices:     make(map[string]schema.Index),
		Constraints: make(map[string]schema.Constraint),
		Sequences:   make(map[string]schema.Sequence),
		MatViews:    make(map[string]schema.MatView),
		Views:       make(map[string]schema.View),
		Triggers:    make(map[string]schema.Trigger),
		Functions:   make(map[string]schema.Function),
		Enums:       make(map[string]schema.Enum),
		Extensions:  make(map[string]schema.Extension),
	}

	// 1. PostgreSQL Version
	var pgVersion string
	err := e.db.QueryRow("SHOW server_version").Scan(&pgVersion)
	if err == nil {
		snapshot.PostgreSQL = pgVersion
	}

	// 2. Tables & Columns
	if err := e.extractTablesAndColumns(snapshot); err != nil {
		return nil, fmt.Errorf("extract tables: %w", err)
	}

	// 3. Indices
	if err := e.extractIndices(snapshot); err != nil {
		return nil, fmt.Errorf("extract indices: %w", err)
	}

	// 4. Constraints
	if err := e.extractConstraints(snapshot); err != nil {
		return nil, fmt.Errorf("extract constraints: %w", err)
	}

	// 5. Sequences
	if err := e.extractSequences(snapshot); err != nil {
		return nil, fmt.Errorf("extract sequences: %w", err)
	}

	// 6. Views and Materialized Views
	if err := e.extractViews(snapshot); err != nil {
		return nil, fmt.Errorf("extract views: %w", err)
	}
	if err := e.extractMatViews(snapshot); err != nil {
		return nil, fmt.Errorf("extract mat views: %w", err)
	}

	// 7. Triggers
	if err := e.extractTriggers(snapshot); err != nil {
		return nil, fmt.Errorf("extract triggers: %w", err)
	}

	// 8. Functions
	if err := e.extractFunctions(snapshot); err != nil {
		// Ignore function parsing errors gracefully
	}

	// 9. Enums
	if err := e.extractEnums(snapshot); err != nil {
		return nil, fmt.Errorf("extract enums: %w", err)
	}

	// 10. Extensions
	if err := e.extractExtensions(snapshot); err != nil {
		return nil, fmt.Errorf("extract extensions: %w", err)
	}

	return snapshot, nil
}

func (e *Extractor) extractTablesAndColumns(snapshot *schema.SchemaSnapshot) error {
	query := `
		SELECT 
			c.table_name, 
			c.column_name, 
			c.data_type, 
			c.is_nullable, 
			c.column_default, 
			c.character_maximum_length, 
			c.numeric_precision, 
			c.numeric_scale, 
			c.collation_name
		FROM information_schema.columns c
		JOIN information_schema.tables t ON c.table_name = t.table_name
		WHERE c.table_schema = 'public' 
		  AND t.table_type = 'BASE TABLE'
		ORDER BY c.table_name, c.ordinal_position
	`

	rows, err := e.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			tableName   string
			colName     string
			dataType    string
			isNullable  string
			colDefault  *string
			charMaxLen  *int
			numPrec     *int
			numScale    *int
			collation   *string
		)

		if err := rows.Scan(&tableName, &colName, &dataType, &isNullable, &colDefault, &charMaxLen, &numPrec, &numScale, &collation); err != nil {
			return err
		}

		table, exists := snapshot.Tables[tableName]
		if !exists {
			table = schema.Table{
				Name:    tableName,
				Columns: make(map[string]schema.Column),
			}
		}

		nullable := isNullable == "YES"

		table.Columns[colName] = schema.Column{
			Name:         colName,
			Type:         dataType,
			Nullable:     nullable,
			Default:      colDefault,
			MaxLen:       charMaxLen,
			NumericPrec:  numPrec,
			NumericScale: numScale,
			Collation:    collation,
		}

		snapshot.Tables[tableName] = table
	}

	return rows.Err()
}
