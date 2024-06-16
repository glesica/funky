package funky

import (
	"testing"
)

func TestApply(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		vals0 := FromVals[int]()
		vals1 := Apply(vals0, double)
		assertProduces(t, vals1)
	})

	t.Run("should handle a single value iterator", func(t *testing.T) {
		vals0 := FromVals(10)
		vals1 := Apply(vals0, double)
		assertProduces(t, vals1, 20)
	})

	t.Run("should handle many iterator values", func(t *testing.T) {
		vals0 := FromVals(10, 20, 30)
		vals1 := Apply(vals0, double)
		assertProduces(t, vals1, 20, 40, 60)
	})
}
