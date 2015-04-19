package index

import (
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

func TestIndexLoading(t *testing.T) {
	Convey("Test index loading", t, func(c C) {
		ind, err := LoadIndex()
		NoError(c, err)
		NotNil(c, ind)
		j, err := ind.JSON()
		NoError(c, err)

		log.Printf("index = %s", j)
	})
}
