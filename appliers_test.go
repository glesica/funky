package funky

import "testing"

func TestMean(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		nextValue := Slice([]float64{}).Apply(Mean[float64]())
		assertProduces(t, nextValue, []iterCase[float64]{
			{0.0, false, nil},
		})
	})

	t.Run("should create a moving average", func(t *testing.T) {
		nextValue := Slice([]float64{1.0, 2.0, 3.0, 0.0}).Apply(Mean[float64]())
		assertProduces(t, nextValue, []iterCase[float64]{
			{1.0, true, nil},
			{1.5, true, nil},
			{2.0, true, nil},
			{1.5, true, nil},
			{0.0, false, nil},
		})
	})
}
