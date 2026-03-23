package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "komparo",
	Short: "Komparo - PostgreSQL Schema Comparison Tool",
	Long: `Komparo compares the structural schema of a PostgreSQL database
against a defined snapshot to prevent schema drift.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
