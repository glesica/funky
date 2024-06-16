package funky

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestMean(t *testing.T) {
	t.Run("should update a moving average", func(t *testing.T) {
		mean := Mean[int]()

		v0, err := mean(0) // 0
		assert.NoError(t, err)
		assert.Equal(t, 0, v0)

		v1, err := mean(1) // 0, 1
		assert.NoError(t, err)
		assert.Equal(t, 0, v1)

		v2, err := mean(2) // 0, 1, 2
		assert.NoError(t, err)
		assert.Equal(t, 1, v2)

		v3, err := mean(9) // 0, 1, 2, 9
		assert.NoError(t, err)
		assert.Equal(t, 3, v3)
	})
}
