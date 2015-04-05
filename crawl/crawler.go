package crawl

import (
	"github.com/xogeny/impact/recorder"
	"log"
)

type Crawler interface {
	Crawl(r recorder.Recorder, logger *log.Logger) error
}
