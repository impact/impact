package config

import (
	"fmt"

	"github.com/impact/impact/crawl"
)

type Settings struct {
	Indices []string
	Sources []crawl.Crawler
	Choices map[string]string
}

func MakeSettings() Settings {
	return Settings{
		Indices: []string{},
		Sources: []crawl.Crawler{},
		Choices: map[string]string{},
	}
}

func (s Settings) String() string {
	return s.List("")
}

func (s Settings) List(prefix string) string {
	ret := ""
	ret = ret + fmt.Sprintf("%sIndex files to read:\n", prefix)
	for _, index := range s.Indices {
		ret = ret + fmt.Sprintf("%s  %s\n", prefix, index)
	}
	ret = ret + fmt.Sprintf("%sSources:\n", prefix)
	for _, source := range s.Sources {
		ret = ret + fmt.Sprintf("%s  %s\n", prefix, source.String())
	}
	if len(s.Choices) == 0 {
		ret = ret + fmt.Sprintf("%sDisambiguations: None\n", prefix)
	} else {
		ret = ret + fmt.Sprintf("%sDisambiguations:\n", prefix)
		for name, uri := range s.Choices {
			ret = ret + fmt.Sprintf("%s  %s -> %s\n", prefix, name, uri)
		}
	}
	return ret
}
