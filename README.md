# Funky

Experiments in functional programming patterns for creating
structured, data-centered workflows.

## Iterator

* TODO: Handle errors consistently when transforming the iterator
* TODO: Decide whether errors are passed through when more = false

The `iter.Func[T]` type represents an iterator of values of type T.
Each call to the iterator will return three values:

1. `value T` - the next value, which might be a zero value
2. `more bool` - false if the iterator has been exhausted
3. `err error` - any that occurred while producing the value

If `more` is true, then either `value` must be the actual next
value in the sequence, or `err` must be non-nil and reflect a
value that ought to have been produced. If `err` is non-nil,
then `value` will be assumed to be meaningless or incomplete,
and implementations should document how they behave in this
regard.

If `more` is false, then `value` will be assumed to be
meaningless and the iterator must not produce any further
values, though it should allow further calls. On subsequent
calls, `more` must also be false. An error can be returned
as well, but its meaning is left undefined.

* `any, false, any` - iterator is exhausted and will produce no
  more values on subsequent calls and the value is meaningless
* `any, true, nil` - the value is valid and there may be
  additional values available through subsequent calls
* `any, true, non-nil` - an error occurred trying to produce
  the next value, which is either meaningless or incomplete,
  depending on the implementation

### Examples

```go
// Create an iterator over the numbers 1-5 from a slice
nums := iter.Slice([]int{1, 2, 3, 4, 5})
```

Then we can use values methods on the `iter.Func[T]` type to
process the values produced by the iterator.

```go
// Sum the even values
s, _ := nums.Where(func(v int) {
    return v % 2 == 0, nil
}).Coalesce(func(a int, v int) (int, error) {
    return a + v, nil
})
```

```go
// Sum the evens and odds separately
nums1, nums2
```
