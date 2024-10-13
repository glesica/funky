package funky

import (
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestNoError(t *testing.T) {
	t.Run("should remove errors from iterator", func(t *testing.T) {
		iter := makeFrom([]Elem[int]{
			{
				val: 1,
				err: nil,
			},
			{
				val: 0,
				err: errors.New("error"),
			},
			{
				val: 2,
				err: nil,
			},
			{
				val: 0,
				err: errors.New("error"),
			},
		})

		iter = NoError(iter)
		assertValues(t, iter, []int{1, 2}, true)
	})

	t.Run("should work when there are only errors", func(t *testing.T) {
		iter := makeFrom([]Elem[int]{
			{
				val: 0,
				err: errors.New("error"),
			},
			{
				val: 0,
				err: errors.New("error"),
			},
		})

		iter = NoError(iter)
		assertValues(t, iter, []int{}, true)
	})

	t.Run("should stop correctly", func(t *testing.T) {
		//
	})
}

func TestTake(t *testing.T) {
	t.Run("should take n elements from an iterator", func(t *testing.T) {
		iter := makeFinite(3)

		iter = Take(iter, 2)
		assertValues(t, iter, []int{0, 1}, true)
	})

	t.Run("should take fewer than n from an exhausted iterator", func(t *testing.T) {
		iter := makeFinite(2)

		iter = Take(iter, 3)
		assertValues(t, iter, []int{0, 1}, true)
	})

	t.Run("should take zero elements from a stopped iterator", func(t *testing.T) {
		iter := makeFinite(0)

		iter = Take(iter, 1)
		assertValues(t, iter, []int{}, true)
	})

	t.Run("should pass through errors", func(t *testing.T) {
		//
	})

	t.Run("should not take after stop", func(t *testing.T) {
		//
	})
}

func TestWhere(t *testing.T) {
	t.Run("should remove elements from the beginning", func(t *testing.T) {
		iter := makeFinite(4)

		iter = Where(iter, func(v int) bool {
			return v > 1
		})
		assertValues(t, iter, []int{2, 3}, true)
	})

	t.Run("should remove elements from the end", func(t *testing.T) {
		iter := makeFinite(4)

		iter = Where(iter, func(v int) bool {
			return v < 2
		})
		assertValues(t, iter, []int{0, 1}, true)
	})

	t.Run("should remove elements from the middle", func(t *testing.T) {
		iter := makeFinite(4)

		iter = Where(iter, func(v int) bool {
			return v != 1 && v != 2
		})
		assertValues(t, iter, []int{0, 3}, true)
	})

	t.Run("should stop correctly", func(t *testing.T) {
		//
	})
}

func TestWhile(t *testing.T) {
	t.Run("should stop delivering values when one fails", func(t *testing.T) {
		iter := makeFinite(10)

		iter = While(iter, func(v int) bool {
			return (v+1)%4 > 0
		})
		assertValues(t, iter, []int{0, 1, 2}, true)
	})

	t.Run("should handle an empty iterator", func(t *testing.T) {
		iter := makeFinite(0)

		iter = While(iter, func(v int) bool {
			return true
		})
		assertValues(t, iter, []int{}, true)
	})
}

func TestZip(t *testing.T) {
	t.Run("should combine values from non-empty iterators", func(t *testing.T) {
		left := makeFinite(2)
		right := makeFinite(2)

		iter := Zip(left, right)
		assertValues(t, iter, []Pair[int, int]{
			{0, 0},
			{1, 1},
		}, true)
	})

	t.Run("should consume but omit extra values", func(t *testing.T) {
		left := makeFinite(2)
		right := makeFinite(1)

		iter := Zip(left, right)
		assertValues(t, iter, []Pair[int, int]{
			{0, 0},
		}, true)

		_, valid := left.Next()
		assert.False(t, valid)
	})

	t.Run("should handle empty iterators", func(t *testing.T) {
		left := makeFinite(0)
		right := makeFinite(0)

		iter := Zip(left, right)
		assertValues(t, iter, []Pair[int, int]{}, true)
	})
}
