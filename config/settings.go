package config

import (
	"fmt"

	"github.com/xogeny/impact/crawl"
)

type Settings struct {
	Indices []string
	Sources []crawl.Crawler
}

func MakeSettings() Settings {
	return Settings{
		Indices: []string{},
		Sources: []crawl.Crawler{},
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
	return ret
}
