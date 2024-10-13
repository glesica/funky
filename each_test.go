package funky

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestEach(t *testing.T) {
	t.Run("should call function for each element", func(t *testing.T) {
		size := 3

		it := makeFinite(size)

		r := &Recorder[int]{}
		it = Each(it, r.Call)
		for range size + 1 {
			_, _ = it.Next()
		}

		assert.Equal(t, []int{0, 1, 2}, r.values)
		assert.Equal(t, []error{nil, nil, nil}, r.errors)
	})
}

type Recorder[T any] struct {
	values []T
	errors []error
}

func (r *Recorder[T]) Call(v T, err error) {
	r.values = append(r.values, v)
	r.errors = append(r.errors, err)
}
