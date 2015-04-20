package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/xogeny/impact/config"
	"github.com/xogeny/impact/index"
)

type IndexCommand struct {
	Output  string `short:"o" long:"output" description:"Output file"`
	Verbose bool   `short:"v" long:"verbose" description:"Turn on verbose output"`
}

func (x IndexCommand) Execute(args []string) error {
	logger := log.New(os.Stdout, "", 0)

	ind := index.NewIndex()

	if x.Output == "" {
		x.Output = "impact_index.json"
	}

	settings, err := config.ReadSettings()
	if err != nil {
		return fmt.Errorf("Error reading settings: %v", err)
	}

	for _, cr := range settings.Sources {
		err = cr.Crawl(ind, x.Verbose, logger)
		if err != nil {
			return fmt.Errorf("Error indexing modelica-3rdparty: %v", err)
		}
	}

	str, _ := ind.JSON()

	if x.Output == "-" {
		fmt.Printf("%s\n", str)
	} else {
		ioutil.WriteFile(x.Output, []byte(str), os.ModePerm)
	}
	return nil
}
