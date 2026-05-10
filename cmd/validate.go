package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rda-run/komparo/internal/config"
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
		dbFlag, _ := cmd.Flags().GetString("db")
		file, _ := cmd.Flags().GetString("file")
		format, _ := cmd.Flags().GetString("format")

		dbURL, err := config.GetConnectionString(dbFlag)
		if err != nil {
			return err
		}

		var data []byte

		if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
			resp, httpErr := http.Get(file)
			if httpErr != nil {
				return fmt.Errorf("failed to fetch remote snapshot file: %w", httpErr)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to fetch remote snapshot file: HTTP %d", resp.StatusCode)
			}
			data, err = io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read remote snapshot file body: %w", err)
			}
		} else {
			data, err = os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("failed to read local snapshot file: %w", err)
			}
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
	validateCmd.Flags().StringP("db", "d", "", "PostgreSQL connecting string (optional if .env is present)")
	validateCmd.Flags().StringP("file", "f", "komparo-snapshot.json", "Input JSON snapshot local file path or remote URL (required)")
	validateCmd.Flags().StringP("format", "o", "text", "Output format (text, json)")
	validateCmd.MarkFlagRequired("file")
}
