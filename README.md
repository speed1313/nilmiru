# nilmiru
![CI](https://img.shields.io/github/actions/workflow/status/speed1313/nilmiru/go.yml?branch=main&label=test)

nilmiru is a static analysis tool that detects nil check leakage in function.

## Features
Points out nil check leakage in function. If the check is leaked, the function may cause panic.

### Example
nilmiru points out nil check leakage.
```go
func slice_invalid(s []string) bool {
	s[0] = "a" // want "nil check leakage"
	return true
}
```
### Constrains
- For now, nilmiru only checks the type of pointer(*hoge) and slice([]hoge), so it can not check composite type field.
- nilmiru forces to do nil check even though the code is valid like below.
```go
func slice_valid_range(s []string) bool{
	for i := range s{ // want "nil check leakage"
		s[i] = "a" // want "nil check leakage"
		return true
	}
	return false
}
```




## How to use
```
$ go install github.com/speed1313/nilmiru/cmd/nilmiru
$ go vet -vettool=$(which nilmiru) ./...
```


