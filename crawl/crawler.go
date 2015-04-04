package crawl

type Crawler interface {
	Crawl(r Recorder) error
}
