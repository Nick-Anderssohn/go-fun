# go-fun
**This is a work in progress**

A pretty self-explanatory streaming/functional programming library. APIs exist for both slices and maps.
Go is not intended to be a functional programming language, so
I'm mainly just creating this for fun and to learn go generics, since they did
not exist last time I used this language.
Slice example:
```golang
	fun.NewSliceStream(testSlice).
		Filter(
			func(v string) (bool, error) {
				return v != "", nil
			},
		).
		Map(
			func(v string) (string, error) {
				return strings.ToUpper(v), nil
			},
		).
		Collect()
```
A similar API exists for maps:
```golang
	fun.NewMapStream(testMap).
		Filter(
			func(k, v string) (bool, error) {
				return k != "" && v != "", nil
			},
		).
		Map(
			func(k, v string) (string, string, error) {
				return strings.ToUpper(k), strings.ToUpper(v), nil
			},
		).
		Collect()
```