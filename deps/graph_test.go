package deps

import "testing"
import "github.com/stretchr/testify/assert"

func TestResolution(t *testing.T) {
	index := MakeLibraryIndex()
	root1 := MakeSemVer(1, 0, 0)
	a1 := MakeSemVer(1, 0, 0)
	index.AddDependency("Root", root1, "A", a1)

	rootVers := index.Versions("Root")
	assert.Equal(t, rootVers, []SemanticVersion{root1})
}
