package funky

import "cmp"

// todo: add sum()

// Mean produces a moving average.
func Mean[T ~int | ~int64 | ~int32 | ~float64 | ~float32]() Applier[T] {
	n := 0
	var total T
	return func(t T) (T, error) {
		n++
		total += t
		return total / T(n), nil
	}
}

// Max produces the maximum value seen so far.
func Max[T cmp.Ordered]() Applier[T] {
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
func Min[T cmp.Ordered]() Applier[T] {
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
