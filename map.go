package funky

// Map returns an iterator that produces the values from the given
// iterator, with the given operation applied to them.
//
// Example (in pseudocode):
//
//	Map({1, 2, 3}, isEven) -> {false, true, false}
func Map[T, R any](it Iter[T], f Mapper[T, R], opts ...Opt) Iter[R] {
	conf := buildConfig(opts...)

	return func(yield func(R) bool) {
		for v := range it {
			r, err := f(v)
			if err != nil {
				conf.Error(err)
				continue
			}

			more := yield(r)
			if !more {
				return
			}
		}
	}
}
