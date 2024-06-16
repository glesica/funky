package funky

import "cmp"

// Mean produces a moving average. If the original iterator
// provides integers, the mean will be computed using integer
// division and will therefore be somewhat inaccurate.
func Mean[T ~int | ~int64 | ~int32 | ~float64 | ~float32]() Mapper[T, T] {
	n := 0
	var total T
	return func(t T) (T, error) {
		n++
		total += t
		return total / T(n), nil
	}
}

// Max produces the maximum value seen so far.
func Max[T cmp.Ordered]() Mapper[T, T] {
	maximum := *new(T)
	first := true

	return func(t T) (T, error) {
		if t > maximum || first {
			maximum = t
			first = false
		}

		return maximum, nil
	}
}

// Min produces the minimum value seen so far.
func Min[T cmp.Ordered]() Mapper[T, T] {
	minimum := *new(T)
	first := true

	return func(t T) (T, error) {
		if t < minimum || first {
			minimum = t
			first = false
		}

		return minimum, nil
	}
}

func Sum[T ~int | ~int64 | ~int32 | ~float64 | ~float32]() Mapper[T, T] {
	var total T
	return func(v T) (T, error) {
		total += v
		return total, nil
	}
}
