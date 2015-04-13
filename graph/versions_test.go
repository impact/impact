package graph

import (
	"testing"

	"github.com/blang/semver"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

func TestVersionConstruction(t *testing.T) {
	Convey("Testing version construction", t, func(c C) {
		vl := NewVersionList()
		Equals(c, 0, vl.Len())

		v1, err := semver.Parse("1.0.0")
		NoError(c, err)

		vl.Add(v1)
		Equals(c, 1, vl.Len())
	})
}

func TestOrdering(t *testing.T) {
	Convey("Testing version construction", t, func(c C) {
		vers := NewVersionList()

		names := []string{
			"1.1.0",
			"0.0.2",
			"1.0.0",
			"1.0.2-beta.1",
			"1.0.2",
			"0.2.1+a8cc9de3f",
			"0.1.1-alpha1+HOTFIX1",
			"0.0.1",
			"2.1.0",
			"1.0.2-alpha.1",
			"1.0.1",
			"0.1.0",
		}

		for _, n := range names {
			v, err := semver.Parse(n)
			NoError(c, err)
			vers.Add(v)
		}

		vers.Sort()

		for i := 0; i < vers.Len()-1; i++ {
			v1 := vers.Get(i)
			v2 := vers.Get(i + 1)
			IsTrue(c, v1.LT(v2))
			IsTrue(c, v2.GT(v1))
			Equals(c, 0, v1.Compare(v1))
			Equals(c, 0, v2.Compare(v2))
		}

		vers.ReverseSort()

		for i := 0; i < vers.Len()-1; i++ {
			v2 := vers.Get(i)
			v1 := vers.Get(i + 1)
			IsTrue(c, v1.LT(v2))
			IsTrue(c, v2.GT(v1))
			Equals(c, 0, v1.Compare(v1))
			Equals(c, 0, v2.Compare(v2))
		}
	})
}
