package deps

import "log"
import "testing"
import "github.com/blang/semver"
import "github.com/stretchr/testify/assert"

func TestResolution(t *testing.T) {
	index := MakeLibraryIndex()
	root1, err := semver.New("1.0.0")
	assert.NoError(t, err, "Parsing root1 version")
	a1, err := semver.New("1.0.0")
	assert.NoError(t, err, "Parsing a1 version")
	index.AddDependency("Root", root1, "A", a1)

	rootVers := index.Versions("Root")
	assert.Equal(t, *rootVers, VersionList{root1})

	// Introduce a circular dependency.  Make sure nothing breaks
	index.AddDependency("A", a1, "Root", root1)

	rootVers = index.Versions("Root")
	assert.Equal(t, *rootVers, VersionList{root1})

	config, err := index.Resolve("Root")
	assert.NoError(t, err)

	log.Printf("Configuration: %v", config)
}
