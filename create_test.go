package funky

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestSlice(t *testing.T) {
	t.Run("should end immediately on a nil slice", func(t *testing.T) {
		nextValue := Slice[int](nil)

		value, more, err := nextValue()
		assert.Equal(t, 0, value)
		assert.False(t, more)
		assert.Zero(t, err)

		_, more, err = nextValue()
		assert.False(t, more)
		assert.Zero(t, err)
	})

	t.Run("should end immediately on an empty slice", func(t *testing.T) {
		nextValue := Slice([]int{})

		value, more, err := nextValue()
		assert.Equal(t, 0, value)
		assert.False(t, more)
		assert.Zero(t, err)

		_, more, err = nextValue()
		assert.False(t, more)
		assert.Zero(t, err)
	})

	t.Run("should produce sequential values", func(t *testing.T) {
		nextValue := Slice([]int{10, 20, 30})

		for _, expectedValue := range []int{10, 20, 30} {
			value, more, err := nextValue()

			assert.Equal(t, expectedValue, value)
			assert.True(t, more)
			assert.Zero(t, err)
		}

		_, more, err := nextValue()
		assert.False(t, more)
		assert.Zero(t, err)
	})
}

func TestEntries(t *testing.T) {
	t.Run("should end immediately on an nil map", func(t *testing.T) {
		nextValue := Entries[int, int](nil)

		value, more, err := nextValue()
		assert.Equal(t, MapEntry[int, int]{}, value)
		assert.False(t, more)
		assert.Zero(t, err)

		_, more, err = nextValue()
		assert.False(t, more)
		assert.Zero(t, err)
	})

	t.Run("should end immediately on an empty map", func(t *testing.T) {
		nextValue := Entries(map[int]int{})

		value, more, err := nextValue()
		assert.Equal(t, MapEntry[int, int]{}, value)
		assert.False(t, more)
		assert.Zero(t, err)

		_, more, err = nextValue()
		assert.False(t, more)
		assert.Zero(t, err)
	})

	t.Run("should produce all values (but not sequentially)", func(t *testing.T) {
		collection := map[int]int{10: 100, 20: 200, 30: 300}
		nextValue := Entries(collection)

		for range []int{0, 1, 2} {
			entry, more, err := nextValue()

			assert.True(t, more)
			assert.Zero(t, err)

			value, ok := collection[entry.Key]
			assert.True(t, ok)
			assert.Equal(t, value, entry.Value)
		}

		_, more, err := nextValue()
		assert.False(t, more)
		assert.Zero(t, err)
	})
}

func TestChan(t *testing.T) {
	t.Run("should produce a valid iterator for a closed channel", func(t *testing.T) {
		c := make(chan int)
		close(c)

		nextValue := Chan(c)
		assertProduces(t, nextValue, []iterCase[int]{
			{0, false, nil},
		})
	})

	t.Run("should produce values from a channel", func(t *testing.T) {
		c := make(chan int, 5)
		for _, i := range []int{10, 20, 30} {
			c <- i
		}
		close(c)

		nextValue := Chan(c)
		assertProduces(t, nextValue, []iterCase[int]{
			{10, true, nil},
			{20, true, nil},
			{30, true, nil},
			{0, false, nil},
		})
	})
}
