package filter

import "regexp"

// Filter collects allowlist and blocklist expressions, and allows callers to
// check if a given string should be permitted. The zero value of a filter type
// is useful and permits all strings.
type Filter struct {
	allowlist []*regexp.Regexp
	blocklist []*regexp.Regexp
}

// Allow adds a regular expression to the allowlist. If the allowlist is
// non-empty, strings must match at least one allowlist expression in order to
// be permitted.
func (f *Filter) Allow(expr string) error {
	re, err := regexp.Compile(expr)
	if err != nil {
		return err
	}

	f.allowlist = append(f.allowlist, re)
	return nil
}

// Block adds a regular expression to the blocklist. If a string matches any
// blocklist expression, it is not permitted.
func (f *Filter) Block(expr string) error {
	re, err := regexp.Compile(expr)
	if err != nil {
		return err
	}

	f.blocklist = append(f.blocklist, re)
	return nil
}

// Permit checks if the provided string is permitted, according to the current
// set of allowlist and blocklist expressions.
func (f *Filter) Permit(s string) (allowed bool) {
	return f.passAllowlist(s) && f.passBlocklist(s)
}

// Blocked checks if the provided string is blocked, according to the current
// blocklist expressions.
func (f *Filter) Blocked(s string) (blocked bool) {
	return !f.passBlocklist(s)
}

func (f *Filter) passAllowlist(s string) bool {
	if len(f.allowlist) <= 0 {
		return true // default pass
	}

	for _, re := range f.allowlist {
		if re.MatchString(s) {
			return true
		}
	}

	return false
}

func (f *Filter) passBlocklist(s string) bool {
	for _, re := range f.blocklist {
		if re.MatchString(s) {
			return false
		}
	}

	return true
}
