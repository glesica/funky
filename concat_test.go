package funky

import "testing"

func TestConcat(t *testing.T) {
	t.Run("should produce nothing when no iterators given", func(t *testing.T) {
		vals0 := Concat[int]()
		assertProduces(t, vals0)
	})

	t.Run("should produce nothing from an empty iterator", func(t *testing.T) {
		vals0 := Concat(FromVals[int]())
		assertProduces(t, vals0)
	})

	t.Run("should reproduce a single iterator", func(t *testing.T) {
		vals0 := Concat(FromVals(1, 2))
		assertProduces(t, vals0, 1, 2)
	})

	t.Run("should concatenate two iterators", func(t *testing.T) {
		vals0 := Concat(FromVals(1, 2), FromVals(3, 4))
		assertProduces(t, vals0, 1, 2, 3, 4)
	})
}
