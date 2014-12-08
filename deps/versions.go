package deps

import "fmt"
import "sort"
import "github.com/blang/semver"

type VersionList []*semver.Version

func NewVersionList(init ...*semver.Version) *VersionList {
	ret := VersionList{}
	ret = append(ret, init...)
	return &ret
}

func (vl *VersionList) Clone() *VersionList {
	return NewVersionList((*vl)...)
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
		if vl2.Contains(v1) {
			*ret = append(*ret, v1)
		}
	}
	return ret
}

func (vl VersionList) Contains(v *semver.Version) bool {
	for _, x := range vl {
		if x.Compare(v) == 0 {
			return true
		}
	}
	return false
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
