package main

import (
	"log"
	"os"
	"testing"

	//"github.com/blang/semver"
	"github.com/xogeny/impact/crawl"
	"github.com/xogeny/impact/index"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

func TestGitHub(t *testing.T) {
	Convey("Testing GitHub crawler", t, func(c C) {
		logger := log.New(os.Stdout, "impact: ", 0)
		cr := crawl.MakeGitHubCrawler("modelica-3rdparty", "")
		ind := index.NewIndex()
		err := cr.Crawl(ind, logger)
		NoError(c, err)

		str, err := ind.JSON()
		NoError(c, err)

		log.Printf("%s", str)
	})
}
