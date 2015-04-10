package parsing

import (
	"fmt"
	"regexp"

	"github.com/blang/semver"
)

// This function takes a version string and returns a semantic version
// representation.  If the string is not, itself, in semantic version
// form, a set of rules will be used to try and cast it into that
// form.
func NormalizeVersion(v string) (semver.Version, error) {
	ret, err := semver.Parse(v)
	if err == nil {
		return ret, nil
	}

	c1, err := regexp.Compile("^[0-9]+\\.[0-9]+")
	if err != nil {
		panic(err)
	}

	m := c1.FindString(v)
	if m != "" {
		nv := fmt.Sprintf("%s.0%s", m, v[len(m):])
		ret, err := semver.Parse(nv)
		if err == nil {
			return ret, nil
		}
	}

	return semver.Version{}, fmt.Errorf("Unable to normalize version string '%s'", v)
}
