package funky

import (
	"errors"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
)

func TestFunc_Apply(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		nextValue := Slice([]int{}).Apply(func(i int) (int, error) {
			return 0, nil
		})
		assertProduces(t, nextValue, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should apply a function to iterator values", func(t *testing.T) {
		nextValue := Slice([]int{10, 20, 30}).Apply(func(i int) (int, error) {
			return i + 1, nil
		})
		assertProduces(t, nextValue, []iterCase[int]{
			{11, true, nil},
			{21, true, nil},
			{31, true, nil},
			{0, false, nil},
		})
	})
}

func TestFunc_Buffered(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		nextValue := Slice([]int{})
		assertProduces(t, nextValue, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should buffer values", func(t *testing.T) {
		actual := []int{10, 20, 30}
		var seen []int

		i := -1
		nextValue := Func[int](func() (int, bool, error) {
			if i >= len(actual) {
				return 0, false, nil
			}

			i++
			seen = append(seen, actual[i])
			return actual[i], true, nil
		}).Buffered(2)

		time.Sleep(1 * time.Second)
		assert.Equal(t, actual[:2], seen)

		value, more, err := nextValue()
		assert.Equal(t, actual[0], value)
		assert.True(t, more)
		assert.Zero(t, err)

		time.Sleep(1 * time.Second)
		assert.Equal(t, actual, seen)
	})
}

func TestFunc_DropN(t *testing.T) {
	t.Run("should skip zero values", func(t *testing.T) {
		actual := []int{10, 20, 30}
		nextValue := Slice(actual).DropN(0)
		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{20, true, nil},
			{30, true, nil},
			{0, false, nil},
		})
	})

	t.Run("should skip one value", func(t *testing.T) {
		actual := []int{10, 20, 30}
		nextValue := Slice(actual).DropN(1)
		assertProduces(t, nextValue, []iterCase[int]{
			{20, true, nil},
			{30, true, nil},
			{0, false, nil},
		})
	})

	t.Run("should skip multiple values", func(t *testing.T) {
		actual := []int{10, 20, 30}
		nextValue := Slice(actual).DropN(2)
		assertProduces(t, nextValue, []iterCase[int]{
			{30, true, nil},
			{0, false, nil},
		})
	})
}

func TestFunc_Drop(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		nextValue := Slice([]int{}).Drop(isEven)
		assertProduces(t, nextValue, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should handle an iterator with an empty prefix", func(t *testing.T) {
		nextValue := Slice([]int{11, 22, 33}).Drop(isEven)
		assertProduces(t, nextValue, []iterCase[int]{
			{11, true, nil},
			{22, true, nil},
			{33, true, nil},
			{0, false, nil},
		})
	})

	t.Run("should handle an iterator with an empty suffix", func(t *testing.T) {
		nextValue := Slice([]int{10, 20, 30}).Drop(isEven)
		assertProduces(t, nextValue, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should handle an iterator with a prefix and suffix", func(t *testing.T) {
		nextValue := Slice([]int{10, 20, 11, 30, 40}).Drop(isEven)
		assertProduces(t, nextValue, []iterCase[int]{
			{11, true, nil},
			{30, true, nil},
			{40, true, nil},
			{0, false, nil},
		})
	})
}

func TestFunc_Parallel(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		nextValue := Slice([]int{}).Parallel(2)
		assertProduces(t, nextValue, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should handle a single-item iterator", func(t *testing.T) {
		nextValue := Slice([]int{10}).Parallel(2)
		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{0, false, nil},
		})
	})

	t.Run("should handle a many-item iterator", func(t *testing.T) {
		nextValue := Slice([]int{10, 20, 30, 40, 50}).Parallel(2)
		values, err := nextValue.Slice()
		assert.Zero(t, err)

		slices.Sort(values)
		assert.Equal(t, []int{10, 20, 30, 40, 50}, values)
	})

	// todo: Test that it actually fans them out correctly
}

func TestFunc_Split(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		nextValue1, nextValue2 := Slice([]int{}).Split()
		assertProduces(t, nextValue1, []iterCase[int]{
			{0, false, nil},
		})
		assertProduces(t, nextValue2, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should produce the same values", func(t *testing.T) {
		nextValue1, nextValue2 := Slice([]int{10, 20}).Split()

		assertProduces(t, nextValue1, []iterCase[int]{
			{10, true, nil},
		})
		assertProduces(t, nextValue2, []iterCase[int]{
			{10, true, nil},
		})

		assertProduces(t, nextValue2, []iterCase[int]{
			{20, true, nil},
		})
		assertProduces(t, nextValue1, []iterCase[int]{
			{20, true, nil},
		})

		assertProduces(t, nextValue1, []iterCase[int]{
			{0, false, nil},
		})
		assertProduces(t, nextValue2, []iterCase[int]{
			{0, false, nil},
		})
	})
}

func TestFunc_TakeN(t *testing.T) {
	t.Run("should take zero values", func(t *testing.T) {
		actual := []int{10, 20, 30}
		nextValue := Slice(actual).TakeN(0)
		assertProduces(t, nextValue, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should take one value", func(t *testing.T) {
		actual := []int{10, 20, 30}
		nextValue := Slice(actual).TakeN(1)
		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{0, false, nil},
		})
	})

	t.Run("should take multiple values", func(t *testing.T) {
		actual := []int{10, 20, 30}
		nextValue := Slice(actual).TakeN(2)
		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{20, true, nil},
			{0, false, nil},
		})
	})
}

func TestFunc_Take(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		nextValue := Slice([]int{}).Take(isEven)
		assertProduces(t, nextValue, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should handle an iterator with an empty prefix", func(t *testing.T) {
		nextValue := Slice([]int{11, 22, 33}).Take(isEven)
		assertProduces(t, nextValue, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should handle an iterator with an empty suffix", func(t *testing.T) {
		nextValue := Slice([]int{10, 20}).Take(isEven)
		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{20, true, nil},
			{0, false, nil},
		})
	})

	t.Run("should handle an iterator with prefix and suffix", func(t *testing.T) {
		nextValue := Slice([]int{10, 20, 11, 30, 40}).Take(isEven)
		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{20, true, nil},
			{0, false, nil},
		})
	})
}

func TestFunc_Where(t *testing.T) {
	t.Run("should handle an exhausted iterator", func(t *testing.T) {
		nextValue := Slice([]int{}).Where(func(i int) (bool, error) {
			return true, nil
		})
		assertProduces(t, nextValue, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should filter first value", func(t *testing.T) {
		nextValue := Slice([]int{10, 20, 30}).Where(func(i int) (bool, error) {
			return i > 10, nil
		})
		assertProduces(t, nextValue, []iterCase[int]{
			{20, true, nil},
			{30, true, nil},
			{0, false, nil},
		})
	})

	t.Run("should filter last value", func(t *testing.T) {
		nextValue := Slice([]int{10, 20, 30}).Where(func(i int) (bool, error) {
			return i < 30, nil
		})
		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{20, true, nil},
			{0, false, nil},
		})
	})

	t.Run("should filter middle value", func(t *testing.T) {
		nextValue := Slice([]int{10, 20, 30}).Where(func(i int) (bool, error) {
			return i != 20, nil
		})
		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{30, true, nil},
			{0, false, nil},
		})
	})
}

func TestFunc_Chan(t *testing.T) {
	t.Run("should handle empty iterator", func(t *testing.T) {
		iterChan := Slice([]int{}).Chan()
		_, more := <-iterChan
		assert.False(t, more)
	})

	t.Run("should produce correct values", func(t *testing.T) {
		actual := []int{10, 20, 30}
		iterChan := Slice(actual).Chan()
		var found []int
		for i := range iterChan {
			found = append(found, i)
		}
		assert.Equal(t, actual, found)
	})
}

func TestFunc_Coalesce(t *testing.T) {
	t.Run("should handle an empty iterator", func(t *testing.T) {
		value, err := Slice([]int{}).Coalesce(sum)
		assert.Zero(t, err)
		assert.Zero(t, value)
	})

	t.Run("should handle a multi-value iterator", func(t *testing.T) {
		value, err := Slice([]int{10, 20, 30}).Coalesce(sum)
		assert.Zero(t, err)
		assert.Equal(t, 60, value)
	})

	// todo: add tests for error conditions
}

func TestFunc_Slice(t *testing.T) {
	t.Run("should handle empty iterator", func(t *testing.T) {
		iterSlice, err := Slice([]int{}).Slice()
		assert.NoError(t, err)
		assert.Equal(t, nil, iterSlice)
	})

	t.Run("should produce correct values", func(t *testing.T) {
		actual := []int{10, 20, 30}
		iterSlice, err := Slice(actual).Slice()
		assert.NoError(t, err)
		assert.Equal(t, actual, iterSlice)
	})

	t.Run("should bubble an error", func(t *testing.T) {
		nextValues := produce([]iterCase[int]{
			{10, true, nil},
			{0, true, errors.New("fake")},
		})
		_, err := nextValues.Slice()
		assert.NotZero(t, err)
	})
}

type iterCase[T any] struct {
	value T
	more  bool
	err   error
}

func assertProduces[T any](t *testing.T, it Func[T], results []iterCase[T]) {
	for i, tc := range results {
		value, more, err := it()
		actual := iterCase[T]{
			value: value,
			more:  more,
			err:   err,
		}

		assert.Equal(t, tc, actual, "element = "+strconv.Itoa(i))
	}
}

func isEven(v int) (bool, error) {
	if v%2 == 0 {
		return true, nil
	}

	return false, nil
}

func sum(acc int, value int) (int, error) {
	return acc + value, nil
}

func produce[T any](iters []iterCase[T]) Func[T] {
	index := 0
	return func() (T, bool, error) {
		if index >= len(iters) {
			return Done[T](nil)
		}

		next := iters[index]
		index++
		return next.value, next.more, next.err
	}
}
