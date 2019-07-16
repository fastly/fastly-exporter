# usage [![builds.sr.ht status](https://builds.sr.ht/~peterbourgon/usage.svg)](https://builds.sr.ht/~peterbourgon/usage?)

Nicer help text for Go programs.

```go
fs := flag.NewFlagSet("my-program", flag.ExitOnError)
var (
	// ...
)
fs.Usage = usage.For(fs, "my-program [flags]")
fs.Parse(os.Args[1:])
```
