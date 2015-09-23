package crawl

import (
	"github.com/impact/impact/recorder"
	"log"
)

type Crawler interface {
	Crawl(r recorder.Recorder, verbose bool, logger *log.Logger) error
	String() string
}
