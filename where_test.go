package funky

import "testing"

func TestWhere(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		vals0 := FromVals[int]()
		vals1 := Where(vals0, func(i int) (bool, error) {
			return true, nil
		})
		assertProduces(t, vals1)
	})

	t.Run("should filter first value", func(t *testing.T) {
		vals0 := FromVals(10, 20, 30)
		vals1 := Where(vals0, func(i int) (bool, error) {
			return i > 10, nil
		})
		assertProduces(t, vals1, 20, 30)
	})

	t.Run("should filter last value", func(t *testing.T) {
		vals0 := FromVals(10, 20, 30)
		vals1 := Where(vals0, func(i int) (bool, error) {
			return i < 30, nil
		})
		assertProduces(t, vals1, 10, 20)
	})

	t.Run("should filter middle value", func(t *testing.T) {
		vals0 := FromVals(10, 20, 30)
		vals1 := Where(vals0, func(i int) (bool, error) {
			return i != 20, nil
		})
		assertProduces(t, vals1, 10, 30)
	})
}
