package crawl

import (
	"log"
	"os"
	"testing"

	"github.com/blang/semver"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"

	"github.com/xogeny/impact/recorder"
)

type NullRecorder struct {
}

func (nr NullRecorder) GetLibrary(owner string, name string) recorder.LibraryRecorder {
	return nr
}

func (nr NullRecorder) SetStars(int)          {}
func (nr NullRecorder) SetEmail(string)       {}
func (nr NullRecorder) SetDescription(string) {}
func (nr NullRecorder) SetHomepage(string)    {}

func (nr NullRecorder) AddVersion(v semver.Version) recorder.VersionRecorder {
	return nr
}

func (nr NullRecorder) SetHash(hash string)                                  {}
func (nr NullRecorder) SetTarballURL(url string)                             {}
func (nr NullRecorder) SetZipballURL(url string)                             {}
func (nr NullRecorder) AddDependency(library string, version semver.Version) {}

func TestGitHub(t *testing.T) {
	Convey("Testing GitHub crawler", t, func(c C) {
		logger := log.New(os.Stdout, "impact: ", 0)
		cr, err := MakeGitHubCrawler("modelica-3rdparty", "", "")
		NoError(c, err)
		err = cr.Crawl(NullRecorder{}, false, logger)
		NoError(c, err)
	})
}

var _ recorder.Recorder = (*NullRecorder)(nil)
var _ recorder.LibraryRecorder = (*NullRecorder)(nil)
var _ recorder.VersionRecorder = (*NullRecorder)(nil)
