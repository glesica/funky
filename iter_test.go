package funky

import (
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestIter_Next(t *testing.T) {
	t.Run("should provide a value", func(t *testing.T) {
		it := &Iter[int]{
			next: func() (Elem[int], bool) {
				return Elem[int]{
					val: 1,
				}, true
			},
		}

		elem, valid := it.Next()
		assert.True(t, valid)
		assert.Equal(t, Elem[int]{val: 1}, elem)
	})
}

func TestIter_Stop(t *testing.T) {
	t.Run("should close the iterator", func(t *testing.T) {
		it := &Iter[int]{
			next: func() (Elem[int], bool) {
				return Elem[int]{
					val: 1,
				}, true
			},
		}

		elem, valid := it.Next()
		assert.True(t, valid)
		assert.Equal(t, Elem[int]{val: 1}, elem)

		it.Close()

		elem, valid = it.Next()
		assert.False(t, valid)
		assert.Equal(t, Elem[int]{val: 0}, elem)
	})
}

func TestIter_ToSeq(t *testing.T) {
	t.Run("should close when the iterator is exhausted", func(t *testing.T) {
		it := makeFinite(3)
		var values []int
		for v := range it.ToSeq() {
			values = append(values, v)
		}

		assert.Equal(t, []int{0, 1, 2}, values)

		elem, valid := it.Next()
		assert.False(t, valid)
		assert.NoError(t, elem.err)
		assert.Equal(t, 0, elem.val)
	})

	t.Run("should close when the loop breaks", func(t *testing.T) {
		it := makeInfinite()
		for v := range it.ToSeq() {
			if v > 0 {
				break
			}
		}

		elem, valid := it.Next()
		assert.True(t, valid)
		assert.NoError(t, elem.err)
		assert.Equal(t, 2, elem.val)
	})
}

func TestFromSlice(t *testing.T) {
	t.Run("should handle an empty slice", func(t *testing.T) {
		it := FromSlice([]int{})
		elem, valid := it.Next()

		assert.False(t, valid)
		assert.Equal(t, Elem[int]{}, elem)
	})

	t.Run("should handle a non-empty slice", func(t *testing.T) {
		vals := []int{1, 2, 3}
		it := FromSlice(vals)

		assertValues(t, it, []int{1, 2, 3}, true)

	})
}

// Creation helpers

func makeFrom(elems []Elem[int]) *Iter[int] {
	i := 0
	return &Iter[int]{
		next: func() (Elem[int], bool) {
			defer func() { i++ }()

			if i >= len(elems) {
				return DoneElem[int]()
			}

			return elems[i], true
		},
		close: func() {
			i = len(elems)
		},
	}
}

func makeInfinite() *Iter[int] {
	n := 0
	return &Iter[int]{
		next: func() (Elem[int], bool) {
			n++
			return ValElem(n - 1)
		},
	}
}

func makeFinite(size int) *Iter[int] {
	n := 0
	return &Iter[int]{
		next: func() (Elem[int], bool) {
			if n >= size {
				return DoneElem[int]()
			}

			n++
			return ValElem(n - 1)
		},
	}
}

func makeErroneous() *Iter[int] {
	return &Iter[int]{
		next: func() (Elem[int], bool) {
			return ErrElem[int](errors.New("error"))
		},
	}
}

func makeConstant[T any](value T) *Iter[T] {
	return &Iter[T]{
		next: func() (Elem[T], bool) {
			return ValElem(value)
		},
	}
}

// Assertion helpers

func assertValues[T any](t *testing.T, it *Iter[T], vals []T, exact bool) {
	for _, value := range vals {
		elem, valid := it.Next()
		assert.True(t, valid)
		assert.NoError(t, elem.err)
		assert.Equal(t, value, elem.val)
	}

	if exact {
		_, valid := it.Next()
		assert.False(t, valid)
	}
}
