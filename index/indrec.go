package index

import (
	"github.com/xogeny/impact/crawl"
)

type IndexRecorder struct {
	index Index
}

func (ir IndexRecorder) AddLibrary(name string) crawl.LibraryRecorder {
	lib := NewLibrary(name)
	ir.index.AddLibrary(lib)
	return makeLibraryRecorder(lib)
}

var _ crawl.Recorder = (*IndexRecorder)(nil)
