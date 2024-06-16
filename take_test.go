package funky

import "testing"

func TestTakeN(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		vals0 := FromVals[int]()
		vals1 := TakeN(vals0, 1)
		assertProduces(t, vals1)
	})

	t.Run("should take zero values", func(t *testing.T) {
		vals0 := FromVals(10, 20, 30)
		vals1 := TakeN(vals0, 0)
		assertProduces(t, vals1)
	})

	t.Run("should take one value", func(t *testing.T) {
		vals0 := FromVals(10, 20, 30)
		vals1 := TakeN(vals0, 1)
		assertProduces(t, vals1, 10)
	})

	t.Run("should take multiple values", func(t *testing.T) {
		vals0 := FromVals(10, 20, 30)
		vals1 := TakeN(vals0, 2)
		assertProduces(t, vals1, 10, 20)
	})
}

func TestTake(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		vals0 := FromVals[int]()
		vals1 := Take(vals0, isEven)
		assertProduces(t, vals1)
	})

	t.Run("should handle an iterator with an empty prefix", func(t *testing.T) {
		vals0 := FromVals(11, 22)
		vals1 := Take(vals0, isEven)
		assertProduces(t, vals1)
	})

	t.Run("should handle an iterator with an empty suffix", func(t *testing.T) {
		vals0 := FromVals(10, 20)
		vals1 := Take(vals0, isEven)
		assertProduces(t, vals1, 10, 20)
	})

	t.Run("should handle an iterator with prefix and suffix", func(t *testing.T) {
		vals0 := FromVals(10, 20, 33, 40)
		vals1 := Take(vals0, isEven)
		assertProduces(t, vals1, 10, 20)
	})
}
