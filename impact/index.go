package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/xogeny/impact/crawl"
	"github.com/xogeny/impact/index"
)

type IndexCommand struct {
	Output  string `short:"o" long:"output" description:"Output file"`
	Verbose bool   `short:"v" long:"verbose" description:"Turn on verbose output"`
}

func (x IndexCommand) Execute(args []string) error {
	logger := log.New(os.Stdout, "", 0)

	ind := index.NewIndex()

	cr, err := crawl.MakeGitHubCrawler("modelica-3rdparty", "", "")
	if err != nil {
		return fmt.Errorf("Error building crawler for modelica-3rdparty: %v", err)
	}
	err = cr.Crawl(ind, false, logger)
	if err != nil {
		return fmt.Errorf("Error indexing modelica-3rdparty: %v", err)
	}

	cr, err = crawl.MakeGitHubCrawler("modelica", "", "")
	if err != nil {
		return fmt.Errorf("Error building crawler for modelica: %v", err)
	}

	err = cr.Crawl(ind, false, logger)
	if err != nil {
		return fmt.Errorf("Error indexing modelica: %v", err)
	}

	str, err := ind.JSON()
	if x.Output != "" {
		ioutil.WriteFile(x.Output, []byte(str), os.ModePerm)
	} else {
		fmt.Printf("%s\n", str)
	}
	return nil
}
