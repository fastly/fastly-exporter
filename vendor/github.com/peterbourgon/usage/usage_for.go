package usage

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
) 

// For returns a usage func for a given flag.FlagSet which outputs to os.Stdout.
// It's sugar for `ForWriter(os.Stdout, fs, short)`.
func For(fs *flag.FlagSet, short string) func() {
	return ForWriter(os.Stdout, fs, short)
}

// ForWriter returns a usage func for a given flag.FlagSet
// which outputs to the given io.Writer.
func ForWriter(w io.Writer , fs *flag.FlagSet, short string) func(){
	return func() {
		fmt.Fprintf(w, "USAGE\n")
		fmt.Fprintf(w, "  %s\n", short)
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "FLAGS\n")
		tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			def := f.DefValue
			if def == "" {
				def = "..."
			}
			fmt.Fprintf(tw, "  -%s %s\t%s\n", f.Name, def, f.Usage)
		})
		tw.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}