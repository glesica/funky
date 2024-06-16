package funky

import (
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		vals0 := FromVals[int]()
		vals1 := Map(vals0, isEven)
		assertProduces(t, vals1)
	})

	t.Run("should handle a single value iterator", func(t *testing.T) {
		vals0 := FromVals(10)
		vals1 := Map(vals0, isEven)
		assertProduces(t, vals1, true)
	})

	t.Run("should handle a multi value iterator", func(t *testing.T) {
		vals0 := FromVals(11, 20, 33)
		vals1 := Map(vals0, isEven)
		assertProduces(t, vals1, false, true, false)
	})
}
