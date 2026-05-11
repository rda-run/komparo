package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rda-run/komparo/internal/config"
	"github.com/rda-run/komparo/internal/db"
	"github.com/rda-run/komparo/internal/migrator"
	"github.com/rda-run/komparo/internal/utils"
	"github.com/rda-run/komparo/pkg/schema"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Generate migration SQL script (read-only operation)",
	Long: `Reads a Komparo JSON snapshot and generates SQL DDL statements
into a file. The user must review and execute the SQL manually.

SECURITY: This command works with read-only database users. No DDL is
ever executed by Komparo.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dbFlag, _ := cmd.Flags().GetString("db")
		file, _ := cmd.Flags().GetString("file")
		output, _ := cmd.Flags().GetString("output")

		dbURL, err := config.GetConnectionString(dbFlag)
		if err != nil {
			return err
		}

		data, err := utils.ReadSnapshot(file)
		if err != nil {
			return err
		}

		var expected schema.SchemaSnapshot
		if err := json.Unmarshal(data, &expected); err != nil {
			return fmt.Errorf("failed to parse snapshot JSON: %w", err)
		}

		extractor, err := db.NewExtractor(dbURL)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		defer extractor.Close()

		actual, err := extractor.ExtractFullSchema()
		if err != nil {
			return fmt.Errorf("failed to extract schema: %w", err)
		}

		changes := migrator.Analyze(&expected, actual)

		if len(changes) == 0 {
			fmt.Println("✅ No differences found. Database already matches the snapshot.")
			return nil
		}

		ordered := migrator.OrderChanges(changes)

		sql := migrator.GenerateFixScript(ordered, dbURL, file)

		outFile := output
		if outFile == "" {
			timestamp := time.Now().Format("20060102_150405")
			outFile = fmt.Sprintf("komparo_fix_%s.sql", timestamp)
		}

		if err := os.WriteFile(outFile, []byte(sql), 0644); err != nil {
			return fmt.Errorf("failed to write fix script: %w", err)
		}

		fmt.Printf("✅ Fix SQL generated: %s\n", outFile)
		fmt.Printf("📊 Found %d differences\n", len(changes))
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Printf("  1. Review the SQL:  less %s\n", outFile)
		fmt.Println("  2. Backup your database:  pg_dump -d <db> > backup.sql")
		fmt.Printf("  3. Execute manually:  psql -d <db> -f %s\n", outFile)
		fmt.Println()
		fmt.Println("⚠️  WARNING: Review the SQL carefully before executing.")
		fmt.Println("   This script may contain DROP statements that delete data.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().StringP("db", "d", "", "PostgreSQL connection string (optional if .env present)")
	migrateCmd.Flags().StringP("file", "f", "", "Snapshot JSON file (required)")
	migrateCmd.Flags().StringP("output", "o", "", "Output SQL file (default: komparo_fix_YYYYMMDD_HHMMSS.sql)")
	migrateCmd.MarkFlagRequired("file")
}
