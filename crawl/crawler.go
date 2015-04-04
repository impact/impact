package crawl

import (
	"log"
)

type Crawler interface {
	Crawl(r Recorder, logger *log.Logger) error
}
