package deps

import "fmt"
import "sort"
import "github.com/blang/semver"

type VersionList []*semver.Version

func NewVersionList() *VersionList {
	return &VersionList{}
}

func (vl *VersionList) Add(v *semver.Version) *VersionList {
	*vl = append(*vl, v)
	return vl
}

func (vl VersionList) Get(i int) semver.Version {
	return *vl[i]
}

func (vl VersionList) Len() int {
	return len(vl)
}

func (vl VersionList) Less(i, j int) bool {
	return vl[i].LT(vl[j])
}

func (vl VersionList) Swap(i, j int) {
	tmp := vl[i]
	vl[i] = vl[j]
	vl[j] = tmp
}

func (vl VersionList) Sort() {
	sort.Sort(vl)
}

func (vl VersionList) ReverseSort() {
	sort.Sort(sort.Reverse(vl))
}

func (vl VersionList) Intersection(vl2 VersionList) *VersionList {
	ret := NewVersionList()
	for _, v1 := range vl {
		for _, v2 := range vl2 {
			if v1.Compare(v2) == 0 {
				*ret = append(*ret, v1)
				break
			}
		}
	}
	return ret
}

func (vl VersionList) String() string {
	ret := "["
	for i, v := range vl {
		if i == 0 {
			ret = fmt.Sprintf("%s%s", ret, (*v).String())
		} else {
			ret = fmt.Sprintf("%s, %s", ret, (*v).String())
		}
	}
	ret = ret + "]"
	return ret
}
