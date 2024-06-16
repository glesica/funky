# Funky

Experiments in functional programming patterns for creating
structured, data-centered workflows.

## `Iter[T]`

The `Iter[T]` type represents an iterator that produces values
of type `T`. In practice, this is just an `iter.Seq[T]`, but it
has some special abilities.

Specifically, it can be used with a variety of functions to create
functional-like data pipelines.

### Examples

```go
// Create an iterator over the numbers 1-5 from a slice
vals := funky.FromSlice([]int{1, 2, 3, 4, 5})
```

Then we can do some transformations:

```go
// Sum the even values
evens := Where(vals, func(v int) (bool, error) {
    return v % 2 == 0, nil
})
total := Coalesce(evens, func(a, v int) (int, error) {
    return a + v, nil
})
```
