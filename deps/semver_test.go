package deps

import "testing"
import "github.com/stretchr/testify/assert"

func TestOrdering(t *testing.T) {
	vers := []SemanticVersion{}

	NewVer := func(major int, minor int, patch int, pre string, post string) {
		v := MakeSemVer(major, minor, patch)
		v.Pre = pre
		v.Post = post
		vers = append(vers, v)
	}

	NewVer(0, 0, 1, "", "")
	NewVer(0, 0, 2, "", "")
	NewVer(0, 1, 0, "", "")
	NewVer(0, 1, 1, "", "")
	NewVer(0, 2, 1, "", "")
	NewVer(1, 0, 0, "", "")
	NewVer(1, 0, 1, "", "")
	NewVer(1, 0, 2, "", "")
	NewVer(1, 1, 0, "", "")
	NewVer(2, 1, 0, "", "")

	for i := 0; i < len(vers)-1; i++ {
		v1 := vers[i]
		v2 := vers[i+1]
		assert.Equal(t, v1.Compare(v2), ComesBefore)
		assert.Equal(t, v2.Compare(v1), ComesAfter)
		assert.Equal(t, v1.Compare(v1), Identical)
		assert.Equal(t, v2.Compare(v2), Identical)
	}
}
