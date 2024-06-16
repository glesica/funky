package funky

// TakeN returns an iterator that produces only the first n values
// from the target.
func TakeN[T any](it Iter[T], n int, opts ...Opt) Iter[T] {
	_ = buildConfig(opts...)

	count := 0
	return func(yield func(T) bool) {
		if n <= 0 {
			return
		}

		for v := range it {
			more := yield(v)
			if !more {
				return
			}

			count++
			if count == n {
				return
			}
		}
	}
}

// Take returns an iterator that produces values from the target
// for as long as the callback returns true, then stops. Therefore,
// calling f on any of the values produced will return true.
//
// Example (in pseudocode):
//
//	Take({1, 1, 2, 1, 1}, x -> x == 1) -> {1, 1}
func Take[T any](it Iter[T], f Predicate[T], opts ...Opt) Iter[T] {
	_ = buildConfig(opts...)

	return func(yield func(T) bool) {
		for v := range it {
			ok, err := f(v)
			if err != nil {
				continue
			}

			if !ok {
				return
			}

			more := yield(v)
			if !more {
				return
			}
		}
	}
}
