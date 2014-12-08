package deps

import "log"
import "testing"
import "github.com/blang/semver"
import "github.com/stretchr/testify/assert"

func TestResolution1(t *testing.T) {
	index := MakeLibraryIndex()
	root1, err := semver.New("1.0.0")
	assert.NoError(t, err, "Parsing root1 version")
	a1, err := semver.New("1.0.0")
	assert.NoError(t, err, "Parsing a1 version")
	index.AddDependency("Root", root1, "A", a1)

	rootVers := index.Versions("Root")
	assert.Equal(t, *rootVers, VersionList{root1})

	config, err := index.Resolve("Root")
	assert.NoError(t, err)
	assert.NotNil(t, config["Root"])

	log.Printf("Configuration: %v", config)

	// Introduce a circular dependency.  Make sure nothing breaks
	index.AddDependency("A", a1, "Root", root1)

	rootVers = index.Versions("Root")
	assert.Equal(t, *rootVers, VersionList{root1})

	config, err = index.Resolve("Root")
	assert.NoError(t, err)

	log.Printf("Configuration: %v", config)
}

func TestResolution2(t *testing.T) {
	index := MakeLibraryIndex()
	root1, err := semver.New("1.0.0")
	assert.NoError(t, err, "Parsing root1 version")
	a1, err := semver.New("1.0.0")
	assert.NoError(t, err, "Parsing a1 version")
	index.AddDependency("Root", root1, "A", a1)

	rootVers := index.Versions("Root")
	assert.Equal(t, *rootVers, VersionList{root1})

	// Introduce a circular dependency that makes resolution fail
	root2, err := semver.New("1.0.1")
	assert.NoError(t, err, "Parsing root2 version")
	index.AddDependency("A", a1, "Root", root2)

	rootVers = index.Versions("Root")
	assert.Equal(t, *rootVers, VersionList{root1})

	config, err := index.Resolve("Root")
	assert.Error(t, err)

	log.Printf("Configuration: %v", config)
}
