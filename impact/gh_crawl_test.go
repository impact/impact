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
		ir := index.IndexRecorder{}
		err := cr.Crawl(ir, logger)
		NoError(c, err)
	})
}
