package funky

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestChunk(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		nextValue := Chunk(Slice([]int{}), 2)

		assertProduces(t, nextValue, []iterCase[[]int]{
			{nil, false, nil},
		})
	})

	t.Run("should chunk a divisible number of values", func(t *testing.T) {
		nextValue := Chunk(Slice([]int{10, 20, 30, 40}), 2)

		assertProduces(t, nextValue, []iterCase[[]int]{
			{[]int{10, 20}, true, nil},
			{[]int{30, 40}, true, nil},
			{nil, false, nil},
		})
	})

	t.Run("should chunk an indivisible number of values", func(t *testing.T) {
		nextValue := Chunk(Slice([]int{10, 20, 30}), 2)

		assertProduces(t, nextValue, []iterCase[[]int]{
			{[]int{10, 20}, true, nil},
			{[]int{30}, true, nil},
			{nil, false, nil},
		})
	})
}

func TestUnchunk(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		nextValue := Unchunk(Slice([][]int{}))

		assertProduces(t, nextValue, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should unchunk homogeneous chunks", func(t *testing.T) {
		nextValue := Unchunk(Slice([][]int{{10, 20}, {30, 40}}))

		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{20, true, nil},
			{30, true, nil},
			{40, true, nil},
			{0, false, nil},
		})
	})

	t.Run("should unchunk heterogeneous chunks", func(t *testing.T) {
		nextValue := Unchunk(Slice([][]int{{10, 20}, {30}}))

		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{20, true, nil},
			{30, true, nil},
			{0, false, nil},
		})
	})
}

func TestConcat(t *testing.T) {
	t.Run("should produce nothing when no iterators given", func(t *testing.T) {
		nextValue := Concat[int]()

		assertProduces(t, nextValue, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should reproduce a single iterator", func(t *testing.T) {
		nextValue := Concat(Slice([]int{10, 20}))

		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{20, true, nil},
			{0, false, nil},
		})
	})

	t.Run("should concatenate two iterators", func(t *testing.T) {
		nextValue := Concat(Slice([]int{10, 20}), Slice([]int{30, 40}))

		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{20, true, nil},
			{30, true, nil},
			{40, true, nil},
			{0, false, nil},
		})
	})
}

func TestEach(t *testing.T) {
	t.Run("should apply to an empty iterator", func(t *testing.T) {
		nextValue := Slice([]int{})

		err := Each(nextValue, func(value int) error {
			return nil
		})
		assert.Zero(t, err)
	})

	t.Run("should apply to a non-empty iterator", func(t *testing.T) {
		nextValue := Slice([]int{10, 20, 30})

		var processedValues []int
		err := Each(nextValue, func(value int) error {
			processedValues = append(processedValues, value*10)
			return nil
		})
		assert.Zero(t, err)
		assert.Equal(t, []int{100, 200, 300}, processedValues)
	})

	t.Run("should stop early on Stop", func(t *testing.T) {
		nextValue := Slice([]int{10, 20, 30})

		var processedValues []int
		err := Each(nextValue, func(value int) error {
			if value > 20 {
				return Stop
			}
			processedValues = append(processedValues, value*10)
			return nil
		})
		assert.Zero(t, err)
		assert.Equal(t, []int{100, 200}, processedValues)
	})
}

func TestMap(t *testing.T) {
	t.Run("should apply function to each value", func(t *testing.T) {
		nextValue := Slice([]int{10, 20, 30})

		mappedValue := Map(nextValue, func(value int) (bool, error) {
			return value == 20, nil
		})

		mappedValues, err := mappedValue.Slice()
		assert.Zero(t, err)
		assert.Equal(t, []bool{false, true, false}, mappedValues)
	})

	// todo: finish map tests

	t.Run("should stop early on Stop", func(t *testing.T) {
		//
	})

	t.Run("should pass through error from callback", func(t *testing.T) {
		//
	})

	t.Run("should pass through error from underlying iterator", func(t *testing.T) {
		//
	})
}
