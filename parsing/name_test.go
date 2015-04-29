package parsing

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

func TestParseName(t *testing.T) {
	Convey("Test name parsing", t, func(c C) {
		name, err := ParseName("package XYZ  blah blah end  XYZ;  ")
		NoError(c, err)
		Equals(c, name, "XYZ")
	})
}
