package funky

import "testing"

func TestZip(t *testing.T) {
	t.Run("should handle empty iterators", func(t *testing.T) {
		left := FromVals[int]()
		right := FromVals[int]()
		both := Zip(left, right)
		assertProduces(t, both)
	})

	t.Run("should handle a shorter left iterator", func(t *testing.T) {
		left := FromVals[int](10)
		right := FromVals[int](20, 30)
		both := Zip(left, right)
		assertProduces(t, both, Pair[int, int]{10, 20}, Pair[int, int]{0, 30})
	})

	t.Run("should handle a shorter right iterator", func(t *testing.T) {
		left := FromVals[int](10, 20)
		right := FromVals[int](30)
		both := Zip(left, right)
		assertProduces(t, both, Pair[int, int]{10, 30}, Pair[int, int]{20, 0})
	})

	t.Run("should handle equal length iterators", func(t *testing.T) {
		left := FromVals[int](10, 20)
		right := FromVals[int](30, 40)
		both := Zip(left, right)
		assertProduces(t, both, Pair[int, int]{10, 30}, Pair[int, int]{20, 40})
	})
}
