package funky

import "cmp"

func Count[I any, A uint64](m A, _ I) (A, error) {
	m++
	return m, nil
}

// todo: create a histogram type of some sort for R

func Histogram[I cmp.Ordered, A any](h A, t I) (A, error) {
	panic("not implemented")
}
