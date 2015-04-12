package parsing

// This function takes a potential complex version string like
// 1.5-build.3 and returns just the version part (no build or
// pre-release strings).  Since we may be dealing with legacy
// version numbers, we cannot assume such a version string is
// a semantic version.  So basically, we cut the string off at the
// first non-numeric non-. character.
func SimpleVersion(v string) string {
	ret := ""
	for _, c := range v {
		if c != '0' && c != '1' && c != '2' && c != '3' && c != '4' && c != '5' &&
			c != '6' && c != '7' && c != '8' && c != '9' && c != '.' {
			break
		}
		ret = ret + string(c)
	}
	return ret
}
