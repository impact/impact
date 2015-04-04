package crawl

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

func TestGitHub(t *testing.T) {
	Convey("Testing GitHub crawler", t, func(c C) {
		cr := MakeGitHubCrawler("mtiller", "")
		err := cr.Crawl()
		NoError(c, err)
	})
}
