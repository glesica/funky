package funky

// A Reducer accepts two values and incorporates the second one
// into the first, which represents some kind of aggregate of
// many values.
//
// For example:
//
//	(sum, value) -> sum + value
type Reducer[I, A any] func(A, I) (A, error)

// Reduce runs a reduce operation on the iterator, calling the
// reducer function repeatedly until the iterator has been exhausted
// or returns an error. If the reducer function returns an error,
// the accumulator will be returned, along with the error.
func Reduce[I any, A any](it *Iter[I], f Reducer[I, A]) (A, error) {
	acc := *new(A)
	var err error

	for elem, valid := it.Next(); valid; elem, valid = it.Next() {
		if elem.err != nil {
			return acc, elem.err
		}

		acc, err = f(acc, elem.val)
		if err != nil {
			return acc, err
		}
	}

	return acc, nil
}
