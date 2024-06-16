package funky

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

// func TestFunc_Buffered(t *testing.T) {
// 	t.Run("should handle an empty iterator", func(t *testing.T) {
// 		nextValue := Slice([]int{})
// 		assertProduces(t, nextValue, []iterCase[int]{
// 			{0, false, nil},
// 		})
// 	})
//
// 	t.Run("should buffer values", func(t *testing.T) {
// 		actual := []int{10, 20, 30}
// 		var seen []int
//
// 		i := -1
// 		nextValue := Iter[int](func() (int, bool, error) {
// 			if i >= len(actual) {
// 				return 0, false, nil
// 			}
//
// 			i++
// 			seen = append(seen, actual[i])
// 			return actual[i], true, nil
// 		}).Buffered(2)
//
// 		time.Sleep(1 * time.Second)
// 		assert.Equal(t, actual[:2], seen)
//
// 		value, more, err := nextValue()
// 		assert.Equal(t, actual[0], value)
// 		assert.True(t, more)
// 		assert.Zero(t, err)
//
// 		time.Sleep(1 * time.Second)
// 		assert.Equal(t, actual, seen)
// 	})
// }
//
// func TestFunc_Parallel(t *testing.T) {
// 	t.Run("should handle an empty iterator", func(t *testing.T) {
// 		nextValue := Slice([]int{}).Parallel(2)
// 		assertProduces(t, nextValue, []iterCase[int]{
// 			{0, false, nil},
// 		})
// 	})
//
// 	t.Run("should handle a single-item iterator", func(t *testing.T) {
// 		nextValue := Slice([]int{10}).Parallel(2)
// 		assertProduces(t, nextValue, []iterCase[int]{
// 			{10, true, nil},
// 			{0, false, nil},
// 		})
// 	})
//
// 	t.Run("should handle a many-item iterator", func(t *testing.T) {
// 		nextValue := Slice([]int{10, 20, 30, 40, 50}).Parallel(2)
// 		values, err := nextValue.Slice()
// 		assert.Zero(t, err)
//
// 		slices.Sort(values)
// 		assert.Equal(t, []int{10, 20, 30, 40, 50}, values)
// 	})
//
// 	// todo: Test that it actually fans them out correctly
// }
//
// func TestFunc_Split(t *testing.T) {
// 	t.Run("should handle an empty iterator", func(t *testing.T) {
// 		nextValue1, nextValue2 := Slice([]int{}).Split()
// 		assertProduces(t, nextValue1, []iterCase[int]{
// 			{0, false, nil},
// 		})
// 		assertProduces(t, nextValue2, []iterCase[int]{
// 			{0, false, nil},
// 		})
// 	})
//
// 	t.Run("should produce the same values", func(t *testing.T) {
// 		nextValue1, nextValue2 := Slice([]int{10, 20}).Split()
//
// 		assertProduces(t, nextValue1, []iterCase[int]{
// 			{10, true, nil},
// 		})
// 		assertProduces(t, nextValue2, []iterCase[int]{
// 			{10, true, nil},
// 		})
//
// 		assertProduces(t, nextValue2, []iterCase[int]{
// 			{20, true, nil},
// 		})
// 		assertProduces(t, nextValue1, []iterCase[int]{
// 			{20, true, nil},
// 		})
//
// 		assertProduces(t, nextValue1, []iterCase[int]{
// 			{0, false, nil},
// 		})
// 		assertProduces(t, nextValue2, []iterCase[int]{
// 			{0, false, nil},
// 		})
// 	})
// }
//
// func TestFunc_Chan(t *testing.T) {
// 	t.Run("should handle empty iterator", func(t *testing.T) {
// 		iterChan := Slice([]int{}).Chan()
// 		_, more := <-iterChan
// 		assert.False(t, more)
// 	})
//
// 	t.Run("should produce correct values", func(t *testing.T) {
// 		actual := []int{10, 20, 30}
// 		iterChan := Slice(actual).Chan()
// 		var found []int
// 		for i := range iterChan {
// 			found = append(found, i)
// 		}
// 		assert.Equal(t, actual, found)
// 	})
// }
//
// func TestFunc_Slice(t *testing.T) {
// 	t.Run("should handle empty iterator", func(t *testing.T) {
// 		iterSlice, err := Slice([]int{}).Slice()
// 		assert.NoError(t, err)
// 		assert.Equal(t, nil, iterSlice)
// 	})
//
// 	t.Run("should produce correct values", func(t *testing.T) {
// 		actual := []int{10, 20, 30}
// 		iterSlice, err := Slice(actual).Slice()
// 		assert.NoError(t, err)
// 		assert.Equal(t, actual, iterSlice)
// 	})
//
// 	t.Run("should bubble an error", func(t *testing.T) {
// 		nextValues := produce([]iterCase[int]{
// 			{10, true, nil},
// 			{0, true, errors.New("fake")},
// 		})
// 		_, err := nextValues.Slice()
// 		assert.NotZero(t, err)
// 	})
// }

// --------------
// Shared helpers
// --------------

func assertProduces[T any](t *testing.T, it Iter[T], expected ...T) {
	var actual []T
	for val := range it {
		actual = append(actual, val)
	}
	assert.Equal(t, expected, actual)
}

func isEven(v int) (bool, error) {
	if v%2 == 0 {
		return true, nil
	}

	return false, nil
}

func double(v int) (int, error) {
	return v * 2, nil
}

func sum(acc int, value int) (int, error) {
	return acc + value, nil
}
