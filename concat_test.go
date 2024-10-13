package funky

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestConcat(t *testing.T) {
	t.Run("should produce nothing when no iterators given", func(t *testing.T) {
		vals0 := Concat[int]()
		assertValues(t, vals0, []int{}, true)
	})

	t.Run("should produce nothing from an empty iterator", func(t *testing.T) {
		vals0 := Concat(FromVals[int]())
		assertValues(t, vals0, []int{}, true)
	})

	t.Run("should reproduce a single iterator", func(t *testing.T) {
		vals0 := Concat(FromVals(1, 2))
		assertValues(t, vals0, []int{1, 2}, true)
	})

	t.Run("should concatenate two iterators", func(t *testing.T) {
		vals0 := Concat(FromVals(1, 2), FromVals(3, 4))
		assertValues(t, vals0, []int{1, 2, 3, 4}, true)
	})

	t.Run("should pass through errors", func(t *testing.T) {
		vals0 := Concat(makeErroneous())
		elem, valid := vals0.Next()
		assert.True(t, valid)
		assert.Error(t, elem.err)
	})
}
