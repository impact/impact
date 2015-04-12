package parsing

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

func TestSimpleVersion(t *testing.T) {
	Convey("Test extracting simple versions", t, func(c C) {
		v1 := SimpleVersion("1.5-build.3")
		Equals(c, v1, "1.5")
	})
}
