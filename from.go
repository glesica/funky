package funky

func FromChan[T any](c <-chan T) Iter[T] {
	return func(yield func(T) bool) {
		for v := range c {
			more := yield(v)
			if !more {
				return
			}
		}
	}
}

func FromMap[K comparable, V any](vals map[K]V) Iter[MapEntry[K, V]] {
	return func(yield func(MapEntry[K, V]) bool) {
		for k, v := range vals {
			more := yield(MapEntry[K, V]{
				Key:   k,
				Value: v,
			})
			if !more {
				return
			}
		}
	}
}

func FromSlice[T any](s []T) Iter[T] {
	return func(yield func(T) bool) {
		for _, v := range s {
			more := yield(v)
			if !more {
				return
			}
		}
	}
}

func FromVals[T any](vs ...T) Iter[T] {
	return FromSlice(vs)
}
