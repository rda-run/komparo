package printer

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rda-run/komparo/internal/diff"
)

const (
	escape = "\xff"
	reset  = escape + "\033[0m" + escape
	red    = escape + "\033[31m" + escape
	yellow = escape + "\033[33m" + escape
	cyan   = escape + "\033[36m" + escape
)

func colorStatus(status string) string {
	switch status {
	case "Missing":
		return red + status + reset
	case "Mismatch":
		return yellow + status + reset
	case "Extra":
		return cyan + status + reset
	default:
		return status
	}
}

func PrintTable(diffs []diff.SchemaDiff) {
	if len(diffs) == 0 {
		fmt.Printf("\n✅ Schema Match: No divergences found.\n\n")
		return
	}
	fmt.Printf("\n⚠️  Schema Divergence Detected!\n\n")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.StripEscape)
	fmt.Fprintln(w, "STATUS\tTYPE\tOBJECT\tDETAILS")

	for _, d := range diffs {
		details := ""
		if d.Expected != "" {
			details += "Expected: " + d.Expected + " "
		}
		if d.Actual != "" {
			if details != "" {
				details += "| "
			}
			details += "Actual: " + d.Actual
		}
		if details == "" {
			details = "<none>"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", colorStatus(d.Status), d.Type, d.Object, details)
	}
	w.Flush()
	fmt.Println()
}

func PrintJSON(diffs []diff.SchemaDiff) {
	data, _ := json.MarshalIndent(diffs, "", "  ")
	fmt.Println(string(data))
}
