package funky

func Chan[T any](c <-chan T) Func[T] {
	return func() (T, bool, error) {
		v, more := <-c
		if !more {
			c = nil
			return v, false, nil
		}

		return v, true, nil
	}
}

func Entries[K comparable, V any](elements map[K]V) Func[MapEntry[K, V]] {
	entries := make([]MapEntry[K, V], len(elements))
	index := 0
	for k, v := range elements {
		entries[index] = MapEntry[K, V]{
			Key:   k,
			Value: v,
		}
		index++
	}

	return Slice(entries)
}

func Keys[K comparable](elements map[K]any) Func[K] {
	keys := make([]K, len(elements))
	index := 0
	for k, _ := range elements {
		keys[index] = k
		index++
	}

	return Slice(keys)
}

func Values[V any](elements map[any]V) Func[V] {
	values := make([]V, len(elements))
	index := 0
	for _, v := range elements {
		values[index] = v
		index++
	}

	return Slice(values)
}

func Slice[T any](elements []T) Func[T] {
	index := -1
	return func() (T, bool, error) {
		index++

		if index >= len(elements) {
			return Done[T](nil)
		}

		return elements[index], true, nil
	}
}
