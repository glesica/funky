package funky

// Where returns an iterator that produces all the values from
// the target for which the given predicate returns true. If the
// predicate returns an error, it will be passed to the OnError
// callback, if provided.
//
// For example (in pseudocode):
//
//	Where({0, 1, 2}, x -> x > 0) -> {1, 2}
func Where[T any](s Iter[T], f Predicate[T], opts ...Opt) Iter[T] {
	params := buildConfig(opts...)

	return func(yield func(T) bool) {
		for v := range s {
			ok, err := f(v)
			if err != nil {
				params.Error(err)
				continue
			}

			if !ok {
				continue
			}

			more := yield(v)
			if !more {
				return
			}
		}
	}
}
