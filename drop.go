package funky

// DropN returns an iterator that produces the same values as the
// target, but skips the first n values.
func DropN[T any](it Iter[T], n int, opts ...Opt) Iter[T] {
	_ = buildConfig(opts...)

	return func(yield func(T) bool) {
		dropped := 0
		for v := range it {
			if dropped < n {
				dropped++
				continue
			}

			more := yield(v)
			if !more {
				return
			}
		}
	}
}

// Drop returns an iterator that skips values for which f returns
// true, then, once it has encountered a value for which f returns
// false, all following values are produced.
//
// For example (in pseudocode):
//
//	Drop({1, 0, 2}, x -> x > 0) -> {0, 2}
func Drop[T any](it Iter[T], f Predicate[T], opts ...Opt) Iter[T] {
	conf := buildConfig(opts...)

	return func(yield func(T) bool) {
		dropped := false
		for v := range it {
			if dropped {
				more := yield(v)
				if !more {
					return
				}

				continue
			}

			ok, err := f(v)
			if err != nil {
				conf.Error(err)
				continue
			}

			if ok {
				continue
			}

			dropped = true

			more := yield(v)
			if !more {
				return
			}
		}
	}
}
