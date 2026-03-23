package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rda-run/komparo/internal/db"
	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Extract schema and save to a JSON snapshot",
	Long:  `Connects to a PostgreSQL database and extracts the schema structure into a JSON file format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dbURL, _ := cmd.Flags().GetString("db")
		outFile, _ := cmd.Flags().GetString("out")
		fmt.Printf("Generating snapshot from %s into %s\n", dbURL, outFile)
		
		extractor, err := db.NewExtractor(dbURL)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		defer extractor.Close()

		snapshot, err := extractor.ExtractFullSchema()
		if err != nil {
			return fmt.Errorf("failed to extract schema: %w", err)
		}

		// Serialize to JSON
		data, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		if err := os.WriteFile(outFile, data, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		fmt.Printf("Successfully generated snapshot to %s\n", outFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.Flags().StringP("db", "d", "", "PostgreSQL connecting string (required)")
	snapshotCmd.Flags().StringP("out", "o", "komparo-snapshot.json", "Output JSON file path")
	snapshotCmd.MarkFlagRequired("db")
}
