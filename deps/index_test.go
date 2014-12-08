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
	index.AddLibrary("Root", root1)
	err = index.AddDependency("Root", root1, "A", a1)
	assert.Error(t, err, "Should fail because A is unknown")
	index.AddLibrary("A", a1)
	err = index.AddDependency("Root", root1, "A", a1)
	assert.NoError(t, err, "Couldn't add dependency")

	rootVers := index.Versions("Root")
	assert.Equal(t, *rootVers, VersionList{root1})

	config, err := index.Resolve("Root")
	assert.NoError(t, err)
	assert.NotNil(t, config["Root"])

	log.Printf("Configuration: %v", config)

	// Introduce a circular dependency.  Make sure nothing breaks
	err = index.AddDependency("A", a1, "Root", root1)
	assert.NoError(t, err, "Couldn't add dependency")

	rootVers = index.Versions("Root")
	assert.True(t, rootVers.Contains(root1))
	assert.Equal(t, 1, rootVers.Len())

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
	index.AddLibrary("Root", root1)
	index.AddLibrary("A", a1)
	err = index.AddDependency("Root", root1, "A", a1)
	assert.NoError(t, err, "Couldn't add dependency")

	rootVers := index.Versions("Root")
	assert.Equal(t, *rootVers, VersionList{root1})

	// Introduce a circular dependency that makes resolution fail
	root2, err := semver.New("1.0.1")
	a2, err := semver.New("1.0.1")
	assert.NoError(t, err, "Parsing root2 version")
	index.AddLibrary("Root", root2)
	index.AddLibrary("A", a2)
	err = index.AddDependency("Root", root2, "A", a2)
	err = index.AddDependency("A", a1, "Root", root2)
	err = index.AddDependency("A", a2, "Root", root1)
	assert.NoError(t, err, "Couldn't add dependency")

	rootVers = index.Versions("Root")
	assert.True(t, rootVers.Contains(root1))
	assert.True(t, rootVers.Contains(root2))
	assert.Equal(t, 2, rootVers.Len())

	config, err := index.Resolve("Root")
	log.Printf("Config = %v", config)
	assert.Error(t, err)

	log.Printf("Configuration: %v", config)
}
