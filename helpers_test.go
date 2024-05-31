package funky

import (
	"errors"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestDone(t *testing.T) {
	t.Run("should return more = false", func(t *testing.T) {
		_, more, _ := Done[any](nil)
		assert.False(t, more)
	})

	t.Run("should return a zero value", func(t *testing.T) {
		value, _, _ := Done[int](nil)
		assert.Zero(t, value)
	})

	t.Run("should return an error when passed one", func(t *testing.T) {
		theErr := errors.New("fake")
		_, _, err := Done[any](theErr)
		assert.Error(t, err)
		assert.Equal(t, theErr, err)
	})
}
