package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rda-run/komparo/internal/db"
	"github.com/rda-run/komparo/internal/diff"
	"github.com/rda-run/komparo/internal/printer"
	"github.com/rda-run/komparo/pkg/schema"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Compare a JSON snapshot against a live database",
	Long:  `Reads a Komparo JSON snapshot and compares it against the structure of a live PostgreSQL database, reporting all differences.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dbURL, _ := cmd.Flags().GetString("db")
		file, _ := cmd.Flags().GetString("file")
		format, _ := cmd.Flags().GetString("format")
		
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read snapshot file: %w", err)
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

		diffs := diff.Compare(&expected, actual)
		
		if format == "json" {
			printer.PrintJSON(diffs)
		} else {
			printer.PrintTable(diffs)
		}

		if len(diffs) > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("db", "d", "", "PostgreSQL connecting string (required)")
	validateCmd.Flags().StringP("file", "f", "komparo-snapshot.json", "Input JSON snapshot file path (required)")
	validateCmd.Flags().StringP("format", "o", "text", "Output format (text, json)")
	validateCmd.MarkFlagRequired("db")
	validateCmd.MarkFlagRequired("file")
}
