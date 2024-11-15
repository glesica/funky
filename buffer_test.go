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
				done <- nil
			}
		})

		b := Buffer(src, uint32(size))
		<-done

		assert.Equal(t, size, count)

		// Consume 1 element and verify that one more is produced
		_, _ = b.Next()
		<-done

		assert.Equal(t, size+1, count)
	})

	t.Run("should stop producing on close", func(t *testing.T) {
		size := 3

		done := make(chan interface{})
		count := 0
		src := Each(makeInfinite(), func(v int, err error) {
			count++
			if count >= size {
				done <- nil
			}
		})

		b := Buffer(src, uint32(size))
		<-done

		// TODO: Add a version that reads two values, spins off a reader that waits until Close has been called, then reads the last

		go assertValues(t, b, []int{0, 1, 2}, true)

		b.Close()

		_, valid := b.Next()
		assert.False(t, valid)
	})
}
