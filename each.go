package funky

// Each calls a function on each value and error produced by
// the input iterator, then passes the element along, unchanged.
// This is intended primarily for logging and the like.
func Each[T any](it *Iter[T], fn func(T, error)) *Iter[T] {
	return &Iter[T]{
		next: func() (Elem[T], bool) {
			elem, valid := it.Next()

			if valid {
				fn(elem.val, elem.err)
			}

			return elem, valid
		},
	}
}
