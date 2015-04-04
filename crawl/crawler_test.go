package crawl

import (
	"log"
	"os"
	"testing"

	"github.com/blang/semver"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xogeny/xconvey"
)

type NullRecorder struct {
}

func (nr NullRecorder) AddLibrary(name string) LibraryRecorder {
	return nr
}

func (nr NullRecorder) SetStars(int) {
}

func (nr NullRecorder) AddVersion(v semver.Version) VersionRecorder {
	return nr
}

func (nr NullRecorder) SetHash(hash string)                                  {}
func (nr NullRecorder) SetTarballURL(url string)                             {}
func (nr NullRecorder) SetZipballURL(url string)                             {}
func (nr NullRecorder) AddDependency(library string, version semver.Version) {}

func TestGitHub(t *testing.T) {
	Convey("Testing GitHub crawler", t, func(c C) {
		logger := log.New(os.Stdout, "impact: ", 0)
		cr := MakeGitHubCrawler("modelica-3rdparty", "")
		err := cr.Crawl(NullRecorder{}, logger)
		NoError(c, err)
	})
}

var _ Recorder = (*NullRecorder)(nil)
var _ LibraryRecorder = (*NullRecorder)(nil)
var _ VersionRecorder = (*NullRecorder)(nil)
