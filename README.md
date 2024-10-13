# Funky

Experiments in functional programming patterns for creating
structured, data-centered workflows.

## `Iter[T]`

The `Iter[T]` type represents an iterator that produces values
of type `T`, along with an error for each value. It can be used
with a variety of functions to create functional-like data pipelines.

### Tools

There are various functions that can be used to transform iterators
to produce different sequences of values. If a transformation
produces an error, then the error will be included. In general,
if a previous step in a pipeline of transformations produced an
error, that error will be passed through untouched. This means that
all errors should survive to the end for handling by the caller.

#### `Apply(...)`

#### `Buffer(...)`

#### `Concat(...)`

#### `Each(...)`

#### `Reduce(...)`

#### `Take(...)`

#### `Where(...)`

#### `Zip(...)`

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
total := Reduce(evens, func(a, v int) (int, error) {
    return a + v, nil
})
```
