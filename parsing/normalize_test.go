package parsing

import (
	"testing"

	"github.com/blang/semver"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

func checkNormalize(c C, s1 string, s2 string) {
	v1, err := NormalizeVersion(s1)
	NoError(c, err)

	v2, err := semver.Parse(s2)
	NoError(c, err)

	Resembles(c, v1, v2)
}

func TestNormalizing(t *testing.T) {
	Convey("Test version normalizing", t, func(c C) {
		checkNormalize(c, "1.2.3", "1.2.3")

		checkNormalize(c, "1.2", "1.2.0")

		checkNormalize(c, "0.5", "0.5.0")

		checkNormalize(c, "0.82", "0.82.0")

		checkNormalize(c, "1.2+build45", "1.2.0+build45")

		_, err := NormalizeVersion("a.b.c")
		IsError(c, err)
	})
}
