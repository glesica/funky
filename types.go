package funky

// An Applier is like a Mapper, but the output is the same type
// as the input.
type Applier[T any] func(T) (T, error)

// A Coalescer is like a reducer, but the accumulator has the same
// type as the inputs. For example, summing a list of numbers.
type Coalescer[T any] func(acc T, value T) (T, error)

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
// an error value.
type Predicate[T any] func(T) (bool, error)

// A Reducer accepts two values and incorporates the second one
// into the first, which represents some kind of aggregate of
// many values.
//
// For example:
//
//	(sum, value) -> sum + value
type Reducer[T any, R any] func(R, T) (R, error)

// errorPair is a helper type that allows us to send full iterator
// values through a channel.
type errorPair[T any] struct {
	value T
	err   error
}
