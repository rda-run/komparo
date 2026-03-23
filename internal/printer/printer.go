package printer

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rda-run/komparo/internal/diff"
)

func PrintTable(diffs []diff.SchemaDiff) {
	if len(diffs) == 0 {
		fmt.Printf("\n✅ Schema Match: No divergences found.\n\n")
		return
	}
	fmt.Printf("\n⚠️  Schema Divergence Detected!\n\n")
	
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
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
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", d.Status, d.Type, d.Object, details)
	}
	w.Flush()
	fmt.Println()
}

func PrintJSON(diffs []diff.SchemaDiff) {
	data, _ := json.MarshalIndent(diffs, "", "  ")
	fmt.Println(string(data))
}
