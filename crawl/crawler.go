package crawl

import (
	"github.com/xogeny/impact/recorder"
	"log"
)

type Crawler interface {
	Crawl(r recorder.Recorder, verbose bool, logger *log.Logger) error
}
