package funky

// Concat creates a single iterator that will produce the values
// from each of the provided iterators, in order.
//
// Note, if an infinite iterator is passed in any but the final
// position, then iterators that follow it will never be used.
//
// For example (in pseudocode):
//
//	Concat({1, 2}, {3, 4}) -> {1, 2, 3, 4}
func Concat[T any](its ...Iter[T]) Iter[T] {
	return func(yield func(T) bool) {
		for _, it := range its {
			for v := range it {
				more := yield(v)
				if !more {
					return
				}
			}
		}
	}
}
