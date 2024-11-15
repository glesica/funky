package funky

// Chunk breaks an iterator into slices of the given size. If there
// are errors, they are joined together and passed along with the
// assembled chunk.
//
// The last chunk may be shorter than size if the number of values
// produced by the underlying iterator is indivisible by size.
//
// This function can be composed with, for example, Each and Buffer,
// to improve the efficiency of parallel computations by performing
// operations on several inputs in one goroutine.
// func Chunk[T any](iterator Iter[T], size int) Iter[[]T] {
// 	return func() ([]T, bool, error) {
// 		next := make([]T, size)
// 		errs := make([]error, size)
//
// 		found := 0
// 		for i := 0; i < size; i++ {
// 			value, more, err := iterator()
// 			if !more {
// 				// There are no additional values, so we're done,
// 				// this means the previous chunk was full and the
// 				// number of values was divisible by size.
// 				if i == 0 {
// 					return Done[[]T](err)
// 				}
//
// 				// The number of values was not divisible by size,
// 				// so we just allow an incomplete chunk.
// 				break
// 			}
//
// 			found++
//
// 			next[i] = value
// 			errs[i] = err
// 		}
//
// 		return Result(next[:found], errors.Join(errs[:found]...))
// 	}
// }

// Unchunk reverses the Chunk operation and properly handles the
// fact that the final chunk may have a different length than
// the others. In fact, all chunks may be arbitrary lengths.
//
// For now, the entire error for a given chunk is passed with
// each of its values. In the future, this function may attempt
// to pair wrapped errors with their corresponding value.
//
// For example:
//
//	Unchunk([]{{1, 2}, {3, 4}}) -> {1, 2, 3, 4}
// func Unchunk[T any](iterator Iter[[]T]) Iter[T] {
// 	var next []T
// 	var errs error
//
// 	return func() (T, bool, error) {
// 		if len(next) == 0 {
// 			value, more, err := iterator()
// 			if !more {
// 				next = nil
// 				errs = nil
//
// 				return Done[T](err)
// 			}
//
// 			next = value
// 			errs = err
// 		}
//
// 		value := next[0]
// 		next = next[1:]
//
// 		return Result(value, errs)
// 	}
// }

// Each applies the given function to each element produced by the
// iterator.
//
// If the function returns Close as its error, iteration will cease
// and no further values will be processed.
//
// Note that this function is not thread-safe unless the iterator
// implements its own synchronization.
//
// For example:
//
//	Each({1, 2, 3}, (v) -> print(v))
// func Each[T any](iterator Iter[T], f func(T) error) error {
// 	for {
// 		value, more, err := iterator()
// 		if err != nil {
// 			return fmt.Errorf("each: iterator error: %w", err)
// 		}
//
// 		if !more {
// 			break
// 		}
//
// 		err = f(value)
// 		if errors.Is(err, Close) {
// 			break
// 		}
// 		if err != nil {
// 			return fmt.Errorf("each: function error: %w", err)
// 		}
// 	}
//
// 	return nil
// }
