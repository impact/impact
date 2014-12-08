package deps

import "fmt"

type OrderRelation int

const ComesBefore OrderRelation = -1
const Identical OrderRelation = 0
const ComesAfter OrderRelation = 1

func (o OrderRelation) String() string {
	if o == ComesBefore {
		return "ComesBefore"
	} else if o == ComesAfter {
		return "ComesAfter"
	} else {
		return "Identical"
	}
}

type SemanticVersion struct {
	Major int
	Minor int
	Patch int
	Pre   string
	Post  string
}

func MakeSemVer(major int, minor int, patch int) SemanticVersion {
	return SemanticVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
		Pre:   "",
		Post:  "",
	}
}

func compInt(v1 int, v2 int) OrderRelation {
	if v1 < v2 {
		return ComesBefore
	} else if v1 > v2 {
		return ComesAfter
	} else {
		return Identical
	}
}

/*
 * This function compares the version 'v1' with
 * the version 'v2'.  If 'v1' comes before 'v2',
 * this function returns -1.  If the versions
 * are equal, it returns 0.  If 'v1' comes after
 * 'v2', this function returns 1.
 */
func (v1 SemanticVersion) Compare(v2 SemanticVersion) OrderRelation {
	major := compInt(v1.Major, v2.Major)
	if major != Identical {
		return major
	}
	minor := compInt(v1.Minor, v2.Minor)
	if minor != Identical {
		return minor
	}
	patch := compInt(v1.Patch, v2.Patch)
	if patch != Identical {
		return patch
	}
	// TODO: Compare pre
	return Identical
}

func (v SemanticVersion) String() string {
	ret := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != "" {
		ret = fmt.Sprintf("%s-%s", ret, v.Pre)
	}
	if v.Post != "" {
		ret = fmt.Sprintf("%s+%s", ret, v.Post)
	}
	return ret
}
