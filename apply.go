package funky

// Apply transforms each input value but keeps the type the same.
// This is a special case of Map where the input and output type
// are identical.
//
// For example:
//
//	Apply({1, 2, 3}, x -> 2*x) -> {2, 4, 6}
func Apply[T any](it Iter[T], f Mapper[T, T], opts ...Opt) Iter[T] {
	conf := buildConfig(opts...)

	return func(yield func(T) bool) {
		for v := range it {
			v, err := f(v)
			if err != nil {
				conf.Error(err)
				continue
			}

			more := yield(v)
			if !more {
				return
			}
		}
	}
}
