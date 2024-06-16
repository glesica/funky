package funky

// A MapEntry simply bundles a key-value pair into a single data
// structure for us with iterators.
type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// A Mapper accepts a value of one type and returns a value of
// a potentially different type, along with an error. It transforms
// values from one set into values from another set.
type Mapper[T any, R any] func(T) (R, error)

// A Pair is just two values, of potentially different types, that
// go together somehow.
type Pair[L any, R any] struct {
	Left  L
	Right R
}

// A Predicate accepts a value and returns true or false, plus
// an error value. If the error is non-nil, the boolean will be
// assumed to be meaningless and ignored.
type Predicate[T any] func(T) (bool, error)

// A Reducer accepts two values and incorporates the second one
// into the first, which represents some kind of aggregate of
// many values.
//
// For example:
//
//	(sum, value) -> sum + value
type Reducer[T any, R any] func(R, T) (R, error)
