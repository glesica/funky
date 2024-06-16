package funky

// Coalesce accumulates the values produced by the target into
// a single value of the same type. Think of it like Reduce, but
// the accumulator has the same type as the inputs.
//
// For example:
//
//	Coalesce({1, 2, 3}, (a, v) -> a + v) -> 6
func Coalesce[T any](it Iter[T], f Reducer[T, T], opts ...Opt) T {
	conf := buildConfig(opts...)

	acc := *new(T)
	for v := range it {
		vv, err := f(acc, v)
		if err != nil {
			conf.Error(err)
			continue
		}

		acc = vv
	}

	return acc
}
