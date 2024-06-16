package funky

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestCoalesce(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		vals0 := FromVals[int]()
		final := Coalesce(vals0, func(acc int, value int) (int, error) {
			return acc + value, nil
		})
		assert.Equal(t, 0, final)
	})

	t.Run("should handle an iterator with one value", func(t *testing.T) {
		vals0 := FromVals(10)
		final := Coalesce(vals0, func(acc int, value int) (int, error) {
			return acc + value, nil
		})
		assert.Equal(t, 10, final)
	})

	t.Run("should handle an iterator with many values", func(t *testing.T) {
		vals0 := FromVals(10, 20, 30)
		final := Coalesce(vals0, func(acc int, value int) (int, error) {
			return acc + value, nil
		})
		assert.Equal(t, 60, final)
	})
}
