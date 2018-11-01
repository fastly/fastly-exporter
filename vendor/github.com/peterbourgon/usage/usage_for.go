package usage

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

// For returns a usage func for a given flag.FlagSet.
func For(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stdout, "USAGE\n")
		fmt.Fprintf(os.Stdout, "  %s\n", short)
		fmt.Fprintf(os.Stdout, "\n")
		fmt.Fprintf(os.Stdout, "FLAGS\n")
		tw := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			def := f.DefValue
			if def == "" {
				def = "..."
			}
			fmt.Fprintf(tw, "  -%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		tw.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
