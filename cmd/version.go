package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is injected during compilation via -ldflags "-X ..."
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Komparo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Komparo version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
