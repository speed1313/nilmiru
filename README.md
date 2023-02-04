# nilmiru

nilmiru is a golang linter which checks nil check leakage in function.

nilmiru Emits a lint error if pointer arguments are used without nil checking.
## Constrains
- For now, nilmiru only checks the type of pointer and slice, so it can not check composite type field.


# Example
```
func slice_invalid(s []string) bool {
	s[0] = "a" // want "nil check leakage"
	return true
}

```