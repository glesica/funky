package funky

import "cmp"

func Count[T any, R uint64](m R, _ T) (R, error) {
	m++
	return m, nil
}

// todo: create a histogram type of some sort for R

func Histogram[T cmp.Ordered, R any](h R, t T) (R, error) {
	panic("not implemented")
}
