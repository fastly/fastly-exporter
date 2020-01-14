package filter

import "regexp"

// Filter collects whitelist and blacklist expressions, and allows callers to
// check if a given string should be permitted. The zero value of a filter type
// is useful and permits all strings.
type Filter struct {
	whitelist []*regexp.Regexp
	blacklist []*regexp.Regexp
}

// Whitelist adds a regular expression to the whitelist. If the whitelist is
// non-empty, strings must match at least one whitelist expression in order to be
// permitted.
func (f *Filter) Whitelist(expr string) error {
	re, err := regexp.Compile(expr)
	if err != nil {
		return err
	}

	f.whitelist = append(f.whitelist, re)
	return nil
}

// Blacklist adds a regular expression to the blacklist. If a string matches any
// blacklist expression, it is not permitted.
func (f *Filter) Blacklist(expr string) error {
	re, err := regexp.Compile(expr)
	if err != nil {
		return err
	}

	f.blacklist = append(f.blacklist, re)
	return nil
}

// Allow checks if the provided string is permitted, according to the current
// set of whitelist and blacklist expressions.
func (f *Filter) Allow(s string) (allowed bool) {
	return f.passWhitelist(s) && f.passBlacklist(s)
}

func (f *Filter) passWhitelist(s string) bool {
	if len(f.whitelist) <= 0 {
		return true // default pass
	}

	for _, re := range f.whitelist {
		if re.MatchString(s) {
			return true
		}
	}

	return false
}

func (f *Filter) passBlacklist(s string) bool {
	for _, re := range f.blacklist {
		if re.MatchString(s) {
			return false
		}
	}

	return true
}
