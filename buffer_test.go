package funky

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestBuffer(t *testing.T) {
	t.Run("should pre-fetch size elements", func(t *testing.T) {
		size := 5

		done := make(chan interface{})
		count := 0
		src := Each(makeInfinite(), func(v int, err error) {
			count++
			if count >= size {
				close(done)
			}
		})

		b := Buffer(src, uint32(size))
		<-done

		assert.Equal(t, size, count)

		// Consume 1 element and verify that one more is produced
		done = make(chan interface{})
		_, _ = b.Next()
		<-done

		assert.Equal(t, size+1, count)
	})
}
