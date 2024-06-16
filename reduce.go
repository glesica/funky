package funky

// Reduce runs a reduce operation on the iterator, calling the
// reducer function repeatedly until the iterator has been exhausted
// and keeping an accumulated value around until it is returned.
func Reduce[T any, R any](it Iter[T], f Reducer[T, R], opts ...Opt) R {
	conf := buildConfig(opts...)

	acc := *new(R)
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
