package deps

import "testing"
import "github.com/blang/semver"
import "github.com/stretchr/testify/assert"

func TestVersionConstruction(t *testing.T) {
	vl := NewVersionList()
	assert.Equal(t, 0, vl.Len())

	v1, err := semver.New("1.0.0")
	assert.NoError(t, err)

	vl.Add(v1)
	assert.Equal(t, 1, vl.Len())
}

func TestOrdering(t *testing.T) {
	vers := NewVersionList()

	names := []string{
		"1.1.0",
		"0.0.2",
		"1.0.0",
		"1.0.2-beta.1",
		"1.0.2",
		"0.2.1+a8cc9de3f",
		"0.1.1-alpha1+HOTFIX1",
		"0.0.1",
		"2.1.0",
		"1.0.2-alpha.1",
		"1.0.1",
		"0.1.0",
	}

	for _, n := range names {
		v, err := semver.New(n)
		assert.NoError(t, err)
		vers.Add(v)
	}

	vers.Sort()

	for i := 0; i < vers.Len()-1; i++ {
		v1 := vers.Get(i)
		v2 := vers.Get(i + 1)
		assert.True(t, (&v1).LT(&v2))
		assert.True(t, (&v2).GT(&v1))
		assert.Equal(t, 0, (&v1).Compare(&v1))
		assert.Equal(t, 0, (&v2).Compare(&v2))
	}

	vers.ReverseSort()

	for i := 0; i < vers.Len()-1; i++ {
		v2 := vers.Get(i)
		v1 := vers.Get(i + 1)
		assert.True(t, (&v1).LT(&v2))
		assert.True(t, (&v2).GT(&v1))
		assert.Equal(t, 0, (&v1).Compare(&v1))
		assert.Equal(t, 0, (&v2).Compare(&v2))
	}
}
