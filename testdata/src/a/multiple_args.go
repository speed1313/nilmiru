package a

func multiple_args(gopher *int, s []string) bool{
	if gopher != nil {
		print(gopher)
	}
	if len(s) != 0{
		s[0] = "a"
		return true
	}
	return false

}

func invalid_multiple_args(gopher *int, s []string) bool{
	print(gopher) // want "nil check leakage"

	s[0] = "a" // want "nil check leakage"
	return true

}