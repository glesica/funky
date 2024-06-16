package funky

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestReduce(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		vals0 := FromVals[int]()
		final := Reduce(vals0, Count[int])
		assert.Equal(t, uint64(0), final)
	})

	t.Run("should handle an iterator with one value", func(t *testing.T) {
		vals0 := FromVals[int](10)
		final := Reduce(vals0, Count[int])
		assert.Equal(t, uint64(1), final)
	})

	t.Run("should handle an iterator with many values", func(t *testing.T) {
		vals0 := FromVals[int](10, 20, 30)
		final := Reduce(vals0, Count[int])
		assert.Equal(t, uint64(3), final)
	})
}
