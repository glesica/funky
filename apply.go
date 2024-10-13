package funky

import "fmt"

// An Applier is a function that can be used with Apply.
type Applier[I, O any] func(I) (O, error)

// Apply transforms each input value, possibly altering the type.
//
// For example (in pseudocode):
//
//	Apply({1, 2, 3}, x -> Letters[x]) -> {"a", "b", "c"}
func Apply[I, O any](it *Iter[I], f Applier[I, O]) *Iter[O] {
	return &Iter[O]{
		next: func() (Elem[O], bool) {
			inElem, valid := it.Next()
			if !valid {
				return DoneElem[O]()
			}

			if inElem.err != nil {
				return ErrElem[O](fmt.Errorf("apply input error: %w", inElem.err))
			}

			outVal, err := f(inElem.val)
			if err != nil {
				return ErrElem[O](err)
			}

			return ValElem(outVal)
		},
	}
}
