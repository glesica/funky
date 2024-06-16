package funky

import "testing"

func TestDropN(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		vals0 := FromVals[int]()
		vals1 := DropN(vals0, 1)
		assertProduces(t, vals1)
	})

	t.Run("should drop zero values", func(t *testing.T) {
		vals0 := FromVals[int](10, 20, 30)
		vals1 := DropN(vals0, 0)
		assertProduces(t, vals1, 10, 20, 30)
	})

	t.Run("should skip one value", func(t *testing.T) {
		vals0 := FromVals[int](10, 20, 30)
		vals1 := DropN(vals0, 1)
		assertProduces(t, vals1, 20, 30)
	})

	t.Run("should skip multiple values", func(t *testing.T) {
		vals0 := FromVals[int](10, 20, 30)
		vals1 := DropN(vals0, 2)
		assertProduces(t, vals1, 30)
	})
}

func TestDrop(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		vals0 := FromVals[int]()
		vals1 := Drop(vals0, isEven)
		assertProduces(t, vals1)
	})

	t.Run("should handle an iterator with an empty prefix", func(t *testing.T) {
		vals0 := FromVals[int](11, 22, 33)
		vals1 := Drop(vals0, isEven)
		assertProduces(t, vals1, 11, 22, 33)
	})

	t.Run("should handle an iterator with an empty suffix", func(t *testing.T) {
		vals0 := FromVals[int](10, 20, 30)
		vals1 := Drop(vals0, isEven)
		assertProduces(t, vals1)
	})

	t.Run("should handle an iterator with a prefix and suffix", func(t *testing.T) {
		vals0 := FromVals[int](10, 20, 33, 40)
		vals1 := Drop(vals0, isEven)
		assertProduces(t, vals1, 33, 40)
	})
}
