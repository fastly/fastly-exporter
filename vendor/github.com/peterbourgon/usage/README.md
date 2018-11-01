# usage

Nicer help text for Go programs.

```go
fs := flag.NewFlagSet("my-program", flag.ExitOnError)
var (
	// ...
)
fs.Usage = usage.For(fs, "my-program [flags]")
fs.Parse(os.Args[1:])
```
