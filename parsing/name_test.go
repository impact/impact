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

		name, err = ParseName(`within ;
package HelmholtzMedia "Data and models of real pure fluids (liquid, two-phase and gas)"
  extends Modelica.Icons.MaterialPropertiesPackage;




  annotation (uses(Modelica(version="3.2.1")));
end HelmholtzMedia;`)
		NoError(c, err)
		Equals(c, name, "HelmholtzMedia")
	})
}
