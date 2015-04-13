package dirinfo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

var sample1 = `
{
  "owner_uri": "https://github.com/sjoelund",
  "libraries": [
          {
                  "name": "MessagePack",
                  "path": "MessagePack"
          }
  ]
}`

func TestDirInfoParsing(t *testing.T) {
	Convey("Test dirinfo parsing", t, func(c C) {
		di, err := Parse(sample1)
		NoError(c, err)
		Equals(c, di.OwnerURI, "https://github.com/sjoelund")
		Equals(c, len(di.Libraries), 1)
		Equals(c, di.Libraries[0].Name, "MessagePack")
		Equals(c, di.Libraries[0].Path, "MessagePack")
		Equals(c, di.Libraries[0].IsFile, false)
	})
}
