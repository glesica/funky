package funky

import (
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestApply(t *testing.T) {
	t.Run("should transform values", func(t *testing.T) {
		it := makeInfinite()
		ait := Apply(it, func(i int) (string, error) {
			return strconv.Itoa(i), nil
		})
		assertValues(t, ait, []string{"0", "1", "2"}, false)
	})

	t.Run("should handle source errors", func(t *testing.T) {
		it := makeErroneous()
		ait := Apply(it, func(i int) (string, error) {
			return strconv.Itoa(i), nil
		})
		elem, valid := ait.Next()
		assert.True(t, valid)
		assert.Equal(t, "", elem.val)
		assert.Error(t, elem.err)
	})

	t.Run("should handle underlying iterator close", func(t *testing.T) {
		it := makeConstant(1)
		ait := Apply(it, func(i int) (string, error) {
			return strconv.Itoa(i), nil
		})

		elem, valid := ait.Next()
		assert.True(t, valid)
		assert.NoError(t, elem.err)
		assert.Equal(t, "1", elem.val)

		it.Close()

		elem, valid = ait.Next()
		assert.False(t, valid)
		assert.NoError(t, elem.err)
		assert.Equal(t, "", elem.val)
	})

	t.Run("should close correctly", func(t *testing.T) {
		it := makeConstant(1)
		ait := Apply(it, func(i int) (string, error) {
			return strconv.Itoa(i), nil
		})

		elem, valid := ait.Next()
		assert.True(t, valid)
		assert.NoError(t, elem.err)
		assert.Equal(t, "1", elem.val)

		ait.Close()

		elem, valid = ait.Next()
		assert.False(t, valid)
		assert.NoError(t, elem.err)
		assert.Equal(t, "", elem.val)

		srcElem, srcValid := it.Next()
		assert.True(t, srcValid)
		assert.NoError(t, srcElem.err)
		assert.Equal(t, 1, srcElem.val)
	})

	t.Run("should handle applier error", func(t *testing.T) {
		//
	})
}
