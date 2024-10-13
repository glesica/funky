package funky

import (
	"cmp"

	"golang.org/x/exp/constraints"
)

// Mean produces a moving average. If the original iterator
// provides integers, the mean will be computed using integer
// division and will therefore be somewhat inaccurate.
func Mean[T Number]() Applier[T, T] {
	n := 0
	var total T
	return func(t T) (T, error) {
		n++
		total += t
		return total / T(n), nil
	}
}

// Max produces the maximum value seen so far.
func Max[T cmp.Ordered]() Applier[T, T] {
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
func Min[T cmp.Ordered]() Applier[T, T] {
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

// Sum produces an iterator that provides a moving sum of the
// values from the underlying iterator.
func Sum[T Number]() Applier[T, T] {
	var total T
	return func(v T) (T, error) {
		total += v
		return total, nil
	}
}

// ToFloat64 converts integer values to float64 values.
func ToFloat64[T constraints.Integer]() Applier[T, float64] {
	return func(v T) (float64, error) {
		return float64(v), nil
	}
}

type Number interface {
	constraints.Integer | constraints.Float
}
