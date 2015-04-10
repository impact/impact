package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/xogeny/impact/crawl"
	"github.com/xogeny/impact/index"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

func TestGitHub(t *testing.T) {
	Convey("Testing GitHub crawler", t, func(c C) {
		logger := log.New(os.Stdout, "impact: ", 0)
		ind := index.NewIndex()

		cr := crawl.MakeGitHubCrawler("modelica-3rdparty", "")
		err := cr.Crawl(ind, false, logger)
		NoError(c, err)

		cr = crawl.MakeGitHubCrawler("modelica", "")
		err = cr.Crawl(ind, false, logger)
		NoError(c, err)

		str, err := ind.JSON()
		NoError(c, err)
		ioutil.WriteFile("gh_crawl.json", []byte(str), os.ModePerm)
	})
}
